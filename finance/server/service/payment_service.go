package service

import (
	"errors"
	"fmt"
	"github.com/akmalfairuz/finance/module/helper"
	"github.com/akmalfairuz/finance/module/random"
	"github.com/akmalfairuz/finance/server/model"
	"github.com/akmalfairuz/finance/server/model/errortype"
	"github.com/akmalfairuz/finance/server/provider"
	"github.com/akmalfairuz/finance/server/repository"
	"github.com/akmalfairuz/finance/server/service/paymentmethod"
	"github.com/go-co-op/gocron/v2"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

type PaymentService struct {
	log                 logrus.FieldLogger
	depositRepository   *repository.DepositRepository
	transactionService  *TransactionService
	notificationService *NotificationService
	depositAlert        provider.TextAlertProvider

	paymentMethods []paymentmethod.PaymentMethod
}

func NewPaymentService(log logrus.FieldLogger, scheduler gocron.Scheduler, depositRepository *repository.DepositRepository) *PaymentService {
	service := &PaymentService{
		log:               log,
		depositRepository: depositRepository,
		paymentMethods:    []paymentmethod.PaymentMethod{},
	}
	if _, err := scheduler.NewJob(gocron.DurationJob(time.Minute), gocron.NewTask(func() {
		if err := service.CheckDepositExpire(); err != nil {
			log.Errorf("Error checking deposit expire: %+v", err)
		}
	}), gocron.WithName("depositExpireChecker"), gocron.WithStartAt(gocron.WithStartImmediately())); err != nil {
		log.Fatalf("Failed to create deposit expire checker scheduler: %+v", err)
	}
	return service
}

func (s *PaymentService) RegisterPaymentMethod(paymentMethod paymentmethod.PaymentMethod) {
	s.paymentMethods = append(s.paymentMethods, paymentMethod)
}

func (s *PaymentService) SetTransactionService(transactionService *TransactionService) {
	s.transactionService = transactionService
}

func (s *PaymentService) SetNotificationService(notificationService *NotificationService) {
	s.notificationService = notificationService
}

func (s *PaymentService) SetDepositAlert(alert provider.TextAlertProvider) {
	s.depositAlert = alert
}

func (s *PaymentService) GetDeposit(depositId int64) (model.DepositRequest, error) {
	return s.depositRepository.GetDeposit(depositId)
}

func (s *PaymentService) CreateDeposit(user model.User, amount int64, paymentMethodId int) (model.DepositRequest, error) {
	waitingDeposit, err := s.depositRepository.GetDepositByUserIdAndStatus(user.ID, model.DepositStatusWaiting)
	if len(waitingDeposit) > 10 {
		return model.DepositRequest{}, errortype.ErrTooManyDeposit
	}

	paymentInfo, err := s.CreatePayment(user.Name, "DEPOSIT", paymentMethodId, amount)
	if err != nil {
		return model.DepositRequest{}, err
	}

	if paymentInfo.Amount-paymentInfo.Fee < 1 {
		return model.DepositRequest{}, fmt.Errorf("error create deposit: invalid total amount %d", paymentInfo.Amount)
	}

	expireIn := time.Duration(paymentInfo.ExpiredAt-time.Now().Unix()) * time.Second

	depositId, err := s.depositRepository.CreateDeposit(model.CreateDepositRequest{
		UserID:              user.ID,
		ExternalPaymentId:   paymentInfo.ExternalPaymentId,
		ExternalPaymentData: paymentInfo.ExternalPaymentData,
		Amount:              paymentInfo.Amount,
		Fee:                 paymentInfo.Fee,
		Description:         fmt.Sprintf("DEPOSIT RP%d VIA %s", paymentInfo.Amount, s.GetPaymentMethodName(paymentMethodId)),
		PaymentMethod:       paymentMethodId,
		ExpireIn:            expireIn,
	})
	if err != nil {
		return model.DepositRequest{}, err
	}

	deposit, err := s.depositRepository.GetDeposit(depositId)
	if err != nil {
		return model.DepositRequest{}, err
	}

	return deposit, nil
}

func (s *PaymentService) getPaymentMethod(paymentMethodId int) (paymentmethod.PaymentMethod, bool) {
	for _, p := range s.paymentMethods {
		if p.ID() == paymentMethodId {
			return p, true
		}
	}
	return nil, false
}

func (s *PaymentService) CreatePayment(customerName string, purpose string, paymentMethodId int, amount int64) (paymentmethod.Data, error) {
	paymentMethod, ok := s.getPaymentMethod(paymentMethodId)
	if !ok {
		return paymentmethod.Data{}, errors.New("payment method not found")
	}

	paymentInfo, err := paymentMethod.CreatePayment(paymentmethod.CreatePaymentOptions{
		RefID:          fmt.Sprintf("%s-%s", purpose, random.String(20)),
		Amount:         amount,
		CustomerName:   customerName,
		ProductDetails: "Pembayaran " + helper.FormatRupiah(amount),
	})
	if err != nil {
		return paymentmethod.Data{}, err
	}

	return paymentInfo, nil
}

func (s *PaymentService) GetDepositPublicInfo(depositId int64) (any, error) {
	deposit, err := s.depositRepository.GetDeposit(depositId)
	if err != nil {
		return nil, err
	}

	paymentMethod, ok := s.getPaymentMethod(deposit.PaymentMethod)
	if !ok {
		return nil, errors.New("payment method not found")
	}

	return paymentMethod.PublicInfo(deposit.ExternalPaymentId, deposit.ExternalPaymentData)
}

func (s *PaymentService) HandlePaymentSuccess(refId string, paymentMethodId int) error {
	paymentPurpose, _, found := strings.Cut(refId, "-")
	if !found {
		return errors.New("invalid ref id")
	}
	switch paymentPurpose {
	case "DEPOSIT":
		depositId, err := s.depositRepository.GetDepositIdByExternalPaymentId(paymentMethodId, refId)
		if err != nil {
			return err
		}
		return s.HandleDepositPaymentSuccess(depositId, paymentMethodId)
	}

	return errors.New(fmt.Sprintf("invalid payment type: %s", paymentPurpose))
}

func (s *PaymentService) HandleDepositPaymentSuccess(depositId int64, paymentMethod int) error {
	deposit, err := s.depositRepository.GetDeposit(depositId)
	if err != nil {
		return err
	}

	if deposit.Status != model.DepositStatusWaiting {
		return errors.New("deposit status is not waiting")
	}

	if deposit.PaymentMethod != paymentMethod {
		return errors.New("payment method is not match")
	}

	tx, err := s.depositRepository.Begin()
	if err != nil {
		return err
	}

	txnId, err := s.transactionService.CreateTransactionWithTx(tx, model.CreateTransaction{
		UserID:      deposit.UserID,
		Type:        model.TransactionTypeDeposit,
		IsAdd:       true,
		ExtraData:   strconv.Itoa(int(deposit.ID)),
		Description: deposit.Description,
		Amount:      deposit.Amount,
	})
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if deposit.Fee > 0 {
		if _, err := s.transactionService.CreateTransactionWithTx(tx, model.CreateTransaction{
			UserID:      deposit.UserID,
			Type:        model.TransactionTypeFee,
			IsAdd:       false,
			ExtraData:   strconv.Itoa(int(deposit.ID)),
			Description: fmt.Sprintf("Biaya admin deposit #%d via %s", deposit.ID, s.GetPaymentMethodName(deposit.PaymentMethod)),
			Amount:      deposit.Fee,
		}); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	if err := s.depositRepository.SetDepositSuccessWithTx(tx, deposit.ID, "Deposit success, balance added to user. Transaction ID: "+strconv.Itoa(int(txnId)), txnId); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		return err
	}

	s.log.WithFields(logrus.Fields{
		"depositId": deposit.ID,
		"userId":    deposit.UserID,
		"amount":    deposit.Amount,
	}).Info("Deposit success")

	if s.depositAlert != nil {
		go func() {
			if err := s.depositAlert.Alert(fmt.Sprintf("[Deposit Success #%d]\nMethod: %s\nAmount: %s\nUserID: %d\n%s", txnId, s.GetPaymentMethodName(deposit.PaymentMethod), helper.FormatRupiah(deposit.Amount), deposit.UserID, helper.CurrentUTC7Date())); err != nil {
				s.log.Errorf("error sending depositAlert: %+v", err)
			}
		}()
	}

	_ = s.notificationService.CreateNotification(model.CreateNotification{
		UserID:      deposit.UserID,
		Title:       fmt.Sprintf("Deposit via %s berhasil", s.GetPaymentMethodName(deposit.PaymentMethod)),
		Description: fmt.Sprintf("Saldo %s telah ditambahkan ke akun anda.", helper.FormatRupiah(deposit.Amount-deposit.Fee)),
		Type:        model.NotificationTypeDeposit,
		Data: map[string]string{
			"_event":    "deposit",
			"depositId": strconv.Itoa(int(depositId)),
		},
		RefID: strconv.Itoa(int(depositId)),
		Push:  true,
	})

	return nil
}

func (s *PaymentService) HandlePaymentFailed(externalId string, paymentMethod int) error {
	deposit, err := s.depositRepository.GetDepositByExternalPaymentId(externalId, paymentMethod)
	if err != nil {
		return err
	}

	tx, err := s.depositRepository.Begin()
	if err != nil {
		return err
	}

	if err := s.depositRepository.SetDepositStatus(deposit.ID, model.DepositStatusFailed, "Deposit failed"); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *PaymentService) GetDepositByUserId(id int64, status *int) ([]model.DepositRequest, error) {
	if status != nil {
		if *status != model.DepositStatusWaiting && *status != model.DepositStatusSuccess && *status != model.DepositStatusFailed {
			return nil, errors.New("invalid status")
		}
		return s.depositRepository.GetDepositByUserIdAndStatus(id, *status)
	}
	return s.depositRepository.GetDepositByUserId(id)
}

func (s *PaymentService) GetPaymentMethodName(paymentMethod int) string {
	paymentMethodObj, ok := s.getPaymentMethod(paymentMethod)
	if !ok {
		return fmt.Sprintf("PAYMENT_%d", paymentMethod)
	}

	return paymentMethodObj.Name()
}

func (s *PaymentService) GetPaymentMethods() []paymentmethod.PaymentMethod {
	ret := make([]paymentmethod.PaymentMethod, 0)
	for _, paymentMethod := range s.paymentMethods {
		ret = append(ret, paymentMethod)
	}
	return ret
}

func (s *PaymentService) CancelDeposit(depositId int64) error {
	deposit, err := s.depositRepository.GetDeposit(depositId)
	if err != nil {
		return err
	}

	if deposit.Status != model.DepositStatusWaiting {
		return errors.New("deposit status is not waiting")
	}

	paymentMethod, ok := s.getPaymentMethod(deposit.PaymentMethod)
	if !ok {
		return fmt.Errorf("unknown payment method %d", deposit.PaymentMethod)
	}

	if err := paymentMethod.Cancel(deposit.ExternalPaymentId); err != nil {
		return err
	}

	return s.depositRepository.SetDepositStatus(deposit.ID, model.DepositStatusFailed, "Deposit canceled by user")
}

func (s *PaymentService) GetDepositTracks(depositId int64) ([]model.DepositRequestTrack, error) {
	return s.depositRepository.GetDepositTracks(depositId)
}

func (s *PaymentService) CheckDepositExpire() error {
	deposits, err := s.depositRepository.GetDepositWillExpire()
	if err != nil {
		return err
	}
	return s.depositRepository.SetDepositExpire(deposits, "Deposit telah kadaluarsa")
}

func (s *PaymentService) GetAllDepositCount(from, to int64) (int64, error) {
	return s.depositRepository.GetAllDepositCount(from, to)
}

func (s *PaymentService) GetTotalDepositAmount(from, to int64) (int64, error) {
	return s.depositRepository.GetTotalDepositAmount(from, to)
}
