package repository

import (
	"github.com/akmalfairuz/finance/module/database"
	"github.com/akmalfairuz/finance/module/hash"
	"github.com/akmalfairuz/finance/server/model"
	"time"
)

type ResetPasswordRepository struct {
	db *database.DB
}

func NewResetPasswordRepository(db *database.DB) *ResetPasswordRepository {
	return &ResetPasswordRepository{db: db}
}

func (r *ResetPasswordRepository) CreateResetPassword(userId int64, token string, expireDuration time.Duration) (int64, error) {
	hashedToken := hash.Sha256Bytes([]byte(token))
	insert, err := r.db.Exec("INSERT INTO resetPasswords (userId, hashedToken, createdAt, expiredAt) VALUES (?, ?, ?, ?)", userId, hashedToken, time.Now().Unix(), time.Now().Add(expireDuration).Unix())
	if err != nil {
		return 0, err
	}
	return insert.LastInsertId()
}

func (r *ResetPasswordRepository) GetResetPassword(id int64) (model.ResetPassword, error) {
	var resetPassword model.ResetPassword
	if err := r.db.Get(&resetPassword, "SELECT * FROM resetPasswords WHERE id = ?", id); err != nil {
		return model.ResetPassword{}, err
	}
	return resetPassword, nil
}

func (r *ResetPasswordRepository) DeleteResetPassword(id int64) error {
	_, err := r.db.Exec("DELETE FROM resetPasswords WHERE id = ?", id)
	return err
}
