package service

import (
	"fmt"
	"github.com/akmalfairuz/finance/module/database"
	"github.com/akmalfairuz/finance/module/helper"
	"github.com/akmalfairuz/finance/server/model"
	"github.com/akmalfairuz/finance/server/repository"
	"github.com/sirupsen/logrus"
	"strconv"
)

type TransferService struct {
	log                 logrus.FieldLogger
	transactionService  *TransactionService
	userService         *UserService
	notificationService *NotificationService
	transferRepository  *repository.TransferRepository
}

func NewTransferService(log logrus.FieldLogger, transferRepository *repository.TransferRepository) *TransferService {
	return &TransferService{log: log, transferRepository: transferRepository}
}

func (s *TransferService) SetTransactionService(transactionService *TransactionService) {
	s.transactionService = transactionService
}

func (s *TransferService) SetUserService(userService *UserService) {
	s.userService = userService
}

func (s *TransferService) SetNotificationService(notificationService *NotificationService) {
	s.notificationService = notificationService
}

func (s *TransferService) GetTransferToUserByUserId(userId int64) ([]model.TransferToUser, error) {
	return s.transferRepository.GetTransferToUser(userId)
}

func (s *TransferService) GetTransferToUserDestinations(userId int64) ([]model.User, error) {
	userIds, err := s.transferRepository.GetTransferToUserDestinations(userId)
	if err != nil {
		return nil, err
	}

	users, err := s.userService.GetUsers(userIds)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *TransferService) TransferToUser(senderUser model.User, receiverUser model.User, amount int64, note string) (int64, error) {
	if amount < 1 {
		return 0, fmt.Errorf("invalid amount")
	}
	tx, err := s.transferRepository.Begin()
	if err != nil {
		return 0, err
	}
	var transferId int64
	err = database.Tx(tx, func() error {
		if _, err := s.userService.GetAndLockUser(tx, senderUser.ID); err != nil {
			return err
		}
		if _, err := s.userService.GetAndLockUser(tx, receiverUser.ID); err != nil {
			return err
		}

		senderTxnID, err := s.transactionService.CreateTransactionWithTx(tx, model.CreateTransaction{
			UserID:      senderUser.ID,
			Type:        model.TransactionTypeTransferToUser,
			IsAdd:       false,
			Description: fmt.Sprintf("Transfer to %s", receiverUser.Email),
			Amount:      amount,
		})
		if err != nil {
			return err
		}
		receiverTxnID, err := s.transactionService.CreateTransactionWithTx(tx, model.CreateTransaction{
			UserID:      receiverUser.ID,
			Type:        model.TransactionTypeReceiveTransferFromUser,
			IsAdd:       true,
			Description: fmt.Sprintf("Received transfer from %s", senderUser.Email),
			Amount:      amount,
		})
		if err != nil {
			return err
		}

		transferId, err = s.transferRepository.CreateTransferToUserWithTx(tx, model.CreateTransferToUser{
			SenderUser:            senderUser,
			ReceiverUser:          receiverUser,
			SenderTransactionID:   senderTxnID,
			ReceiverTransactionID: receiverTxnID,
			Amount:                amount,
			Note:                  note,
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	_ = s.notificationService.CreateNotification(model.CreateNotification{
		UserID:      receiverUser.ID,
		Title:       fmt.Sprintf("Anda menerima %s", helper.FormatRupiah(amount)),
		Description: fmt.Sprintf("%s telah mengirim %s kepada Anda.\nCatatan dari pengirim: %s", senderUser.Name, helper.FormatRupiah(amount), note),
		Type:        model.NotificationTypeTransfer,
		Data: map[string]string{
			"_event":     "refreshBalance",
			"transferId": fmt.Sprintf("%d", transferId),
		},
		RefID: strconv.Itoa(int(transferId)),
		Push:  true,
	})

	return transferId, err
}

func (s *TransferService) GetTransferToUserByID(id int64) (model.TransferToUser, error) {
	return s.transferRepository.GetTransferToUserByID(id)
}
