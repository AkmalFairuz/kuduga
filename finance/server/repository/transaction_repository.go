package repository

import (
	"github.com/akmalfairuz/finance/module/database"
	"github.com/akmalfairuz/finance/server/model"
	"github.com/jmoiron/sqlx"
	"time"
)

type TransactionRepository struct {
	db *database.DB
}

func NewTransactionRepository(db *database.DB) *TransactionRepository {
	return &TransactionRepository{db}
}

func (r *TransactionRepository) CreateTransactionWithTx(tx *sqlx.Tx, userId int64, transactionType int, extraData string, description string, beforeBalance int64, afterBalance int64, amount int64) (int64, error) {
	insert, err := tx.Exec("INSERT INTO transactions (userId, type, extraData, description, amount, beforeBalance, afterBalance, createdAt) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", userId, transactionType, extraData, description, amount, beforeBalance, afterBalance, time.Now().Unix())
	if err != nil {
		return 0, err
	}
	return insert.LastInsertId()
}

func (r *TransactionRepository) GetTransaction(id int64) (model.Transaction, error) {
	var transaction model.Transaction
	if err := r.db.Get(&transaction, "SELECT * FROM transactions WHERE id = ?", id); err != nil {
		return model.Transaction{}, err
	}
	return transaction, nil
}

func (r *TransactionRepository) GetUserTransactions(userId int64, startTime int64, endTime int64, limit int) ([]model.Transaction, error) {
	var transactions []model.Transaction
	if err := r.db.Select(&transactions, "SELECT * FROM transactions WHERE userId = ? AND createdAt >= ? AND createdAt <= ? ORDER BY id DESC LIMIT ?", userId, startTime, endTime, limit); err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *TransactionRepository) GetUserTransactionsBeforeID(userId int64, beforeId int64, limit int) ([]model.Transaction, error) {
	var transactions []model.Transaction
	if err := r.db.Select(&transactions, "SELECT * FROM transactions WHERE userId = ? AND id < ? ORDER BY id DESC LIMIT ?", userId, beforeId, limit); err != nil {
		return nil, err
	}
	return transactions, nil
}
