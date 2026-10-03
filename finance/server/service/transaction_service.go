package service

import (
	"errors"
	"github.com/akmalfairuz/finance/server/model"
	"github.com/akmalfairuz/finance/server/model/errortype"
	"github.com/akmalfairuz/finance/server/repository"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type TransactionService struct {
	log                   logrus.FieldLogger
	transactionRepository *repository.TransactionRepository
	userService           *UserService
}

var ErrTransactionBalanceNotEnough = errortype.ErrBalanceNotEnough

func NewTransactionService(log logrus.FieldLogger, transactionRepository *repository.TransactionRepository) *TransactionService {
	return &TransactionService{
		log:                   log,
		transactionRepository: transactionRepository,
	}
}

func (s *TransactionService) SetUserService(userService *UserService) {
	s.userService = userService
}

func (s *TransactionService) CreateTransactionWithTx(tx *sqlx.Tx, create model.CreateTransaction) (int64, error) {
	if create.Amount <= 0 {
		return 0, errors.New("invalid amount")
	}
	user, err := s.userService.GetAndLockUser(tx, create.UserID)
	if err != nil {
		return 0, err
	}
	if user.Locked {
		return 0, ErrTransactionBalanceNotEnough
	}
	amount := create.Amount
	if !create.IsAdd {
		if user.Balance < create.Amount {
			return 0, ErrTransactionBalanceNotEnough
		}
		amount = -amount
	}

	txnId, err := s.transactionRepository.CreateTransactionWithTx(tx, create.UserID, create.Type, create.ExtraData, create.Description, user.Balance, user.Balance+amount, amount)
	if err != nil {
		return 0, err
	}
	if err := s.userService.AddUserBalanceWithTx(tx, create.UserID, amount); err != nil {
		if errors.Is(err, repository.ErrBalanceNotEnough) {
			return 0, ErrTransactionBalanceNotEnough
		}
		return 0, err
	}
	return txnId, nil
}

func (s *TransactionService) GetUserTransactions(userId int64, startTime int64, endTime int64, limit int) ([]model.Transaction, error) {
	// make start time in 00:00:00 of the day
	startTime = startTime - (startTime % 86400)
	// and end time in 23:59:59 of the day
	endTime = endTime - (endTime % 86400) + 86399
	if startTime > endTime {
		return nil, errors.New("invalid start time and end time")
	}
	return s.transactionRepository.GetUserTransactions(userId, startTime, endTime, limit)
}
