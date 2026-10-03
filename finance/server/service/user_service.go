package service

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"github.com/akmalfairuz/finance/module/database"
	"github.com/akmalfairuz/finance/module/hash"
	"github.com/akmalfairuz/finance/module/helper"
	"github.com/akmalfairuz/finance/module/password"
	"github.com/akmalfairuz/finance/module/pointer"
	"github.com/akmalfairuz/finance/module/random"
	"github.com/akmalfairuz/finance/module/templating"
	"github.com/akmalfairuz/finance/server/model"
	"github.com/akmalfairuz/finance/server/model/errortype"
	"github.com/akmalfairuz/finance/server/provider"
	"github.com/akmalfairuz/finance/server/repository"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"html/template"
	"strings"
	"time"
)

type UserService struct {
	log                     logrus.FieldLogger
	emailService            *EmailService
	userRepository          *repository.UserRepository
	resetPasswordRepository *repository.ResetPasswordRepository
	alert                   provider.TextAlertProvider

	shutdownEmailTemplate *template.Template
}

//go:embed email_shutdown_template.html
var emailClosingTemplate string

func NewUserService(log logrus.FieldLogger, userRepository *repository.UserRepository, resetPasswordRepository *repository.ResetPasswordRepository) *UserService {
	tmpl, err := templating.LoadMinified("email.html", emailClosingTemplate)
	if err != nil {
		log.Fatalf("failed loading email shutdown template", err)
	}
	return &UserService{
		log:                     log,
		userRepository:          userRepository,
		resetPasswordRepository: resetPasswordRepository,
		shutdownEmailTemplate:   tmpl,
	}
}

func (s *UserService) SetEmailService(emailService *EmailService) {
	s.emailService = emailService
}

func (s *UserService) SetAlert(alert provider.TextAlertProvider) {
	s.alert = alert
}

func (s *UserService) CreateUser(user model.CreateUser) (model.User, error) {
	hashedPassword, err := password.Hash([]byte(user.Password))
	if err != nil {
		return model.User{}, err
	}
	userId, err := s.userRepository.CreateUser(user.Name, user.DisplayName, user.Email, user.DeviceUniqueID, hashedPassword)
	if err != nil {
		return model.User{}, err
	}
	u, err := s.userRepository.GetUser(userId)
	if err != nil {
		return model.User{}, err
	}
	_ = s.alert.Alert(fmt.Sprintf("[New User #%d]\nUsername: %s\nDisplay Name: %s\nEmail: %s", u.ID, u.DisplayName, u.Email))
	return u, nil
}

func (s *UserService) GetUser(id int64) (model.User, error) {
	return s.userRepository.GetUser(id)
}

func (s *UserService) GetUserByName(name string) (model.User, error) {
	return s.userRepository.GetUserByName(name)
}

func (s *UserService) GetUserByEmail(email string) (model.User, error) {
	return s.userRepository.GetUserByEmail(email)
}

func (s *UserService) CreateResetPassword(userId int64, token string) (model.ResetPassword, error) {
	id, err := s.resetPasswordRepository.CreateResetPassword(userId, token, time.Minute*15)
	if err != nil {
		return model.ResetPassword{}, err
	}
	resetPassword, err := s.resetPasswordRepository.GetResetPassword(id)
	if err != nil {
		return model.ResetPassword{}, err
	}
	return resetPassword, nil
}

func (s *UserService) ResetPassword(id int64, token, newPassword string) error {
	resetPassword, err := s.resetPasswordRepository.GetResetPassword(id)
	if err != nil {
		return err
	}
	if !bytes.Equal(hash.Sha256Bytes([]byte(token)), resetPassword.HashedToken) {
		return errors.New("invalid reset password token")
	}
	if time.Now().Unix() >= resetPassword.ExpiredAt {
		return errors.New("reset password already expired")
	}
	if err := s.resetPasswordRepository.DeleteResetPassword(id); err != nil {
		return err
	}
	newPasswordHash, err := password.Hash([]byte(newPassword))
	if err != nil {
		return err
	}
	return s.userRepository.SetUserPasswordHash(resetPassword.UserID, newPasswordHash)
}

func (s *UserService) AddUserBalanceWithTx(tx *sqlx.Tx, userId int64, amount int64) error {
	return s.userRepository.AddUserBalanceWithTx(tx, userId, amount)
}

func (s *UserService) UpdateUserEmail(userId int64, email string) error {
	return s.userRepository.UpdateUserEmail(userId, email)
}

func (s *UserService) UpdateUserDisplayName(userId int64, displayName string) error {
	return s.userRepository.UpdateUserDisplayName(userId, displayName)
}

func (s *UserService) VerifyPassword(userId int64, pwd string) error {
	user, err := s.userRepository.GetUser(userId)
	if err != nil {
		return err
	}
	return password.Verify(user.PasswordHash, []byte(pwd))
}

func (s *UserService) UpdateUserPassword(userId int64, pwd string) error {
	hashedPassword, err := password.Hash([]byte(pwd))
	if err != nil {
		return err
	}
	return s.userRepository.SetUserPasswordHash(userId, hashedPassword)
}

func (s *UserService) CreateUserPin(userId int64, pin string) error {
	user, err := s.userRepository.GetUser(userId)
	if err != nil {
		return err
	}
	if user.PinHash != nil {
		return errors.New("cannot create pin")
	}
	pinHash, err := password.Hash([]byte(pin))
	if err != nil {
		return err
	}
	if err := s.userRepository.UpdateUserPinHash(userId, &pinHash); err != nil {
		return err
	}
	return nil
}

func (s *UserService) UpdateUserPin(userId int64, pin string, newPin string) error {
	user, err := s.userRepository.GetUser(userId)
	if err != nil {
		return err
	}
	if user.PinHash != nil && password.Verify([]byte(*user.PinHash), []byte(pin)) != nil {
		return errortype.ErrInvalidPin
	}
	newPinHash, err := password.Hash([]byte(newPin))
	if err != nil {
		return err
	}
	if err := s.userRepository.UpdateUserPinHash(userId, &newPinHash); err != nil {
		return err
	}
	return nil
}

func (s *UserService) DeleteUserPin(userId int64, pin string, force bool) error {
	user, err := s.userRepository.GetUser(userId)
	if err != nil {
		return err
	}
	if !force {
		if user.PinHash != nil && password.Verify([]byte(*user.PinHash), []byte(pin)) != nil {
			return errortype.ErrInvalidPin
		}
	}
	if err := s.userRepository.UpdateUserPinHash(userId, nil); err != nil {
		return err
	}
	return nil
}

func (s *UserService) GetUsers(ids []int64) ([]model.User, error) {
	return s.userRepository.GetUsers(ids)
}

func (s *UserService) GetAndLockUser(tx *sqlx.Tx, id int64) (model.User, error) {
	return s.userRepository.GetAndLockUser(tx, id)
}

func (s *UserService) UpdateUserKycStatus(userId int64, status bool) error {
	return s.userRepository.UpdateUserKycStatus(userId, status)
}

func (s *UserService) GenerateUsernameFromGoogleAuth(info *model.GoogleUserInfo) (string, error) {
	// attempt 1, try to use email as username (without domain)
	{
		emailUser, _, found := strings.Cut(info.Email, "@")
		if found {
			// remove prohibited characters
			emailUser = strings.ToLower(emailUser)
			emailUser = helper.StringRemoveNonAlphanumeric(emailUser)

			// get first 18 characters
			if len(emailUser) > 16 {
				emailUser = emailUser[:16]
			}

			if model.IsValidUsername(emailUser) {
				// check if username already exists
				attempt := 0
				for attempt < 5 {
					attempt++
					generatedUsername := emailUser
					if attempt > 1 {
						generatedUsername = fmt.Sprintf("%s%d", emailUser, random.Number(10, 99))
					}
					_, err := s.userRepository.GetUserByName(generatedUsername)
					if errors.Is(err, database.ErrNoRows) {
						return generatedUsername, nil
					}
					if err != nil {
						return "", err
					}
				}
			}
		}
	}

	// attempt 2, try to use display name as username
	{
		displayName := helper.StringRemoveNonAlphanumeric(info.DisplayName)
		displayName = strings.ToLower(displayName)
		if len(displayName) > 16 {
			displayName = displayName[:16]
		}
		if model.IsValidUsername(displayName) {
			// check if username already exists
			attempt := 0
			for attempt < 5 {
				attempt++
				generatedUsername := displayName
				if attempt > 1 {
					generatedUsername = fmt.Sprintf("%s%d", displayName, random.Number(10, 99))
				}
				_, err := s.userRepository.GetUserByName(generatedUsername)
				if errors.Is(err, database.ErrNoRows) {
					return generatedUsername, nil
				}
				if err != nil {
					return "", err
				}
			}
		}
	}

	return fmt.Sprintf("user%d", random.Number(100000, 999999)), nil
}

func (s *UserService) GetNewUserCount(from, to int64) (int64, error) {
	return s.userRepository.GetNewUserCount(from, to)
}

func (s *UserService) GetTotalUserBalance() (int64, error) {
	return s.userRepository.GetTotalUserBalance()
}

func (s *UserService) GenerateReferralCode(user *model.User) (string, error) {
	if user.ReferralCode != nil {
		return *user.ReferralCode, nil
	}
	referralCode := random.StringUpper(5)
	if err := s.userRepository.SetUserReferralCode(user.ID, referralCode); err != nil {
		return "", err
	}
	user.ReferralCode = &referralCode
	return referralCode, nil
}

func (s *UserService) UseReferralCode(user *model.User, code string) error {
	if user.ReferralCode != nil {
		return errors.New("user already has referral code")
	}
	if user.CreatedAt < time.Now().Add(time.Hour*-1).Unix() {
		return errors.New("user cannot use referral code")
	}
	toUser, err := s.userRepository.GetUserByReferralCode(code)
	if user.DeviceUniqueID == nil || toUser.DeviceUniqueID == nil {
		return errortype.ErrReferralUnable
	}
	if *toUser.DeviceUniqueID == *user.DeviceUniqueID || toUser.ID == user.ID {
		return errortype.ErrReferralSelf
	}
	usersCount, err := s.userRepository.GetUserCountByDeviceUniqueIDLast7Days(*user.DeviceUniqueID)
	if usersCount > 1 {
		return errortype.ErrReferralUnable
	}
	if err != nil {
		if errors.Is(err, database.ErrNoRows) {
			return errortype.ErrReferralNotFound
		}
		return err
	}
	user.ReferralCode = &code
	user.ReferralAt = pointer.Make(time.Now().Unix())
	user.ReferralUserID = &toUser.ID
	return s.userRepository.SetUserReferralTo(user.ID, toUser.ID)
}

func (s *UserService) GetCountUsedThisReferralUser(referralUserId int64) (int64, error) {
	return s.userRepository.GetCountUsedThisReferralUser(referralUserId)
}

func (s *UserService) UpdateUserDeviceUniqueID(id int64, deviceUniqueId string) error {
	return s.userRepository.UpdateUserDeviceUniqueID(id, deviceUniqueId)
}

func (s *UserService) AnnounceClosingToUser() error {
	users, err := s.userRepository.GetUsersWithBalance()
	if err != nil {
		return err
	}
	for _, user := range users {
		body := bytes.NewBuffer(nil)

		if err := s.shutdownEmailTemplate.Execute(body, map[string]any{"username": user.DisplayName, "balance": helper.FormatRupiah(user.Balance)}); err != nil {
			return fmt.Errorf("error parsing template: %w", err)
		}

		if err := s.emailService.SendEmail([]string{user.Email}, "[Pemberitahuan Penting #2] Penutupan Layanan Aplikasi Kuduga", body.String(), "Kuduga No-Reply <noreply@example.invalid>"); err != nil {
			s.log.Errorf("error sending closing email to %s: %v", user.Email, err)
		}
		time.Sleep(time.Second)

	}
	return nil
}
