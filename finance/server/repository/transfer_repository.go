package repository

import (
	"github.com/akmalfairuz/finance/module/database"
	"github.com/akmalfairuz/finance/server/model"
	"github.com/jmoiron/sqlx"
	"time"
)

type TransferRepository struct {
	db *database.DB
}

func NewTransferRepository(db *database.DB) *TransferRepository {
	return &TransferRepository{db: db}
}

func (r *TransferRepository) Begin() (*sqlx.Tx, error) {
	return r.db.Beginx()
}

func (r *TransferRepository) CreateTransferToUserWithTx(tx *sqlx.Tx, create model.CreateTransferToUser) (int64, error) {
	result, err := tx.Exec("INSERT INTO transferToUser (senderUserId, senderUsername, senderName, senderEmail, senderTransactionId, receiverUserId, receiverUsername, receiverName, receiverEmail, receiverTransactionId, amount, purpose, note, createdAt) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		create.SenderUser.ID, create.SenderUser.Name, create.SenderUser.DisplayName, create.SenderUser.Email, create.SenderTransactionID, create.ReceiverUser.ID, create.ReceiverUser.Name, create.ReceiverUser.DisplayName, create.ReceiverUser.Email, create.ReceiverTransactionID, create.Amount, create.Purpose, create.Note, time.Now().Unix())
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *TransferRepository) GetTransferToUser(senderUserId int64) ([]model.TransferToUser, error) {
	var ret []model.TransferToUser
	if err := r.db.Select(&ret, "SELECT * FROM transferToUser WHERE senderUserId = ? ORDER BY id DESC", senderUserId); err != nil {
		return nil, err
	}
	return ret, nil
}

func (r *TransferRepository) GetTransferToUserDestinations(senderUserId int64) ([]int64, error) {
	var ret []struct {
		ReceiverUserID int64 `db:"receiverUserId"`
	}
	if err := r.db.Select(&ret, "SELECT receiverUserId FROM transferToUser WHERE senderUserId = ? GROUP BY receiverUserId, id ORDER BY id DESC", senderUserId); err != nil {
		return nil, err
	}
	ids := make([]int64, 0)
	for _, id := range ret {
		ids = append(ids, id.ReceiverUserID)
	}
	return ids, nil
}

func (r *TransferRepository) GetTransferToUserByID(id int64) (model.TransferToUser, error) {
	var ret model.TransferToUser
	if err := r.db.Get(&ret, "SELECT * FROM transferToUser WHERE id = ?", id); err != nil {
		return model.TransferToUser{}, err
	}
	return ret, nil
}
