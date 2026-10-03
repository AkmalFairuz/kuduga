package service

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/akmalfairuz/finance/module/hash"
	"github.com/akmalfairuz/finance/module/random"
	"github.com/akmalfairuz/finance/module/templating"
	"github.com/akmalfairuz/finance/server/model"
	"github.com/akmalfairuz/finance/server/repository"
	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"html/template"
	"time"
)

type OTPService struct {
	log          logrus.FieldLogger
	repository   *repository.OTPRepository
	emailService *EmailService

	rateLimiter *redis_rate.Limiter

	otpEmailFrom string

	otpEmailTemplate *template.Template
}

var ErrOTPInvalidToken = errors.New("invalid otp token")
var ErrOTPInvalidCode = errors.New("invalid otp code")
var ErrOTPInvalidScope = errors.New("invalid otp scope")
var ErrOTPAlreadyUsed = errors.New("otp already used")
var ErrOTPExpired = errors.New("otp expired")
var ErrOTPTooManyFailAttempt = errors.New("too many fail attempt")
var ErrOTPTooManyRetryAttempt = errors.New("too many retry attempt")

var otpScopeDetails = map[int]otpScopeDetail{
	model.OTPScopeRegister: {
		Subject:     "Verifikasi Akun",
		Description: "Anda baru saja melakukan registrasi akun.<br/>Mohon abaikan email ini jika Anda tidak melakukan registrasi.",
	},
	model.OTPScopeResetPassword: {
		Subject:     "Permintaan Reset Kata Sandi",
		Description: "Kami menerima permintaan untuk mengatur ulang kata sandi akun Anda.<br/>Jika Anda tidak melakukan permintaan ini, mohon mengabaikan email ini.",
	},
	model.OTPScopeUpdateEmail: {
		Subject:     "Permintaan Perubahan Alamat Email",
		Description: "Kami menerima permintaan untuk mengubah alamat email akun Anda.<br/>Jika Anda tidak melakukan permintaan ini, mohon untuk mengabaikan email ini.",
	},
	model.OTPScopeRecoveryPin: {
		Subject:     "Permintaan Reset PIN",
		Description: "Kami menerima permintaan untuk mengatur ulang PIN akun Anda.<br/>Jika Anda tidak melakukan permintaan ini, mohon untuk mengabaikan email ini.",
	},
}

type otpScopeDetail struct {
	Subject     string
	Description string
}

//go:embed email_otp_template.html
var emailOtpTemplateContent string

func NewOTPService(log logrus.FieldLogger, otpEmailFrom string, redisClient *redis.Client, r *repository.OTPRepository) *OTPService {
	otpEmailTemplate, err := templating.LoadMinified("email.html", emailOtpTemplateContent)
	if err != nil {
		log.Fatalf("failed to parse otp email template: %v", err)
	}

	if otpEmailFrom == "" {
		log.Fatalf("otp email from is empty")
	}

	return &OTPService{
		log:              log,
		repository:       r,
		otpEmailTemplate: otpEmailTemplate,
		otpEmailFrom:     otpEmailFrom,
		rateLimiter:      redis_rate.NewLimiter(redisClient),
	}
}

func (s *OTPService) SetEmailService(emailService *EmailService) {
	s.emailService = emailService
}

func (s *OTPService) RequestEmailOTP(ipAddress, email string, scope int, rawTok string) (model.OTP, error) {
	if _, err := s.rateLimiter.Allow(context.TODO(), "otp:email:"+ipAddress, redis_rate.PerHour(10)); err != nil {
		return model.OTP{}, err
	}
	code := s.generateCode()

	id, err := s.repository.CreateOTP(model.OTPTypeEmail, scope, email, code, rawTok, time.Minute*30)
	if err != nil {
		return model.OTP{}, err
	}

	otp, err := s.repository.GetOTP(id)
	if err != nil {
		return model.OTP{}, err
	}

	s.sendOTP(otp.Scope, otp.Type, otp.Contact, code)

	return otp, nil
}

func (s *OTPService) sendEmailOTP(email, code string, scope int) error {
	s.log.Infof("sending otp to email %s with code %s", email, code)
	body := bytes.NewBuffer(nil)

	scopeInfo, ok := otpScopeDetails[scope]
	if !ok {
		return fmt.Errorf("invalid otp scope %d", scope)
	}

	if err := s.otpEmailTemplate.Execute(body, map[string]any{"description": template.HTML(scopeInfo.Description), "email": email, "code": code}); err != nil {
		return fmt.Errorf("failed to execute otp email template: %w", err)
	}
	if err := s.emailService.SendEmail([]string{email}, scopeInfo.Subject, body.String(), s.otpEmailFrom); err != nil {
		return fmt.Errorf("failed to send otp email to %s: %w", email, err)
	}
	return nil
}

func (s *OTPService) ResendOTP(id int64, token string) error {
	otp, err := s.repository.GetOTP(id)
	if err != nil {
		return err
	}
	if !bytes.Equal(hash.Sha256Bytes([]byte(token)), otp.HashedToken) {
		return ErrOTPInvalidToken
	}
	if err := s.isInvalidOTP(otp); err != nil {
		return err
	}
	if otp.RetryAttempt >= 5 {
		return ErrOTPTooManyRetryAttempt
	}
	newCode := s.generateCode()
	if err := s.repository.Retry(otp, newCode); err != nil {
		return err
	}

	s.sendOTP(otp.Scope, otp.Type, otp.Contact, newCode)

	return nil
}

func (s *OTPService) sendOTP(scope, otpType int, contact, code string) {
	switch otpType {
	case model.OTPTypeEmail:
		go func() {
			if err := s.sendEmailOTP(contact, code, scope); err != nil {
				s.log.Error(err)
			}
		}()
		return
	}
}

func (s *OTPService) isInvalidOTP(otp model.OTP) error {
	if otp.UsedAt != nil {
		return ErrOTPAlreadyUsed
	}
	if time.Now().Unix() >= otp.ExpiredAt {
		return ErrOTPExpired
	}
	return nil
}

func (s *OTPService) VerifyOTP(id int64, contact string, scope int, token string, code string) error {
	otp, err := s.repository.GetOTP(id)
	if err != nil {
		return err
	}
	if otp.Contact != contact {
		return ErrOTPInvalidToken //TODO: change to ErrOTPInvalidContact
	}
	if !bytes.Equal(hash.Sha256Bytes([]byte(token)), otp.HashedToken) {
		return ErrOTPInvalidToken
	}
	if otp.Scope != scope {
		return ErrOTPInvalidScope
	}
	if otp.FailAttempt >= 5 {
		return ErrOTPTooManyFailAttempt
	}
	if err := s.isInvalidOTP(otp); err != nil {
		return err
	}
	if otp.Code != code {
		if err := s.repository.AddOTPFailAttempt(id); err != nil {
			return err
		}
		return ErrOTPInvalidCode
	}
	return s.repository.UseOTP(otp)
}

func (s *OTPService) GetOTP(id int64) (model.OTP, error) {
	return s.repository.GetOTP(id)
}

func (s *OTPService) generateCode() string {
	const chars = "1234567890"
	return random.StringWithChars(6, chars)
}
