package repository

import (
	"errors"
	"github.com/akmalfairuz/finance/module/database"
	"github.com/akmalfairuz/finance/server/model"
	"time"
)

type UserKycRepository struct {
	db *database.DB
}

func NewUserKycRepository(db *database.DB) *UserKycRepository {
	return &UserKycRepository{db: db}
}

var ErrUserKycWaiting = errors.New("wait")

func (r *UserKycRepository) Create(create *model.CreateUserKyc) (int64, error) {
	var c struct {
		Count int64 `db:"c"`
	}
	if err := r.db.Get(&c, "SELECT COUNT(*) AS c FROM userKyc WHERE userId = ? AND status = 0", create.UserID); err != nil {
		return 0, err
	}
	if c.Count != 0 {
		return 0, ErrUserKycWaiting
	}
	insert, err := r.db.Exec("INSERT INTO userKyc (userId, documentId, fullName, documentFileId, createdAt) VALUES (?, ?, ?, ?, ?)", create.UserID, create.DocumentID, create.FullName, create.DocumentFileID, time.Now().Unix())
	if err != nil {
		return 0, err
	}
	return insert.LastInsertId()
}

func (r *UserKycRepository) UpdateStatus(id int64, status int) error {
	_, err := r.db.Exec("UPDATE userKyc SET status = ? WHERE id = ?", status, id)
	return err
}

func (r *UserKycRepository) GetByUserId(userId int64) (model.UserKyc, error) {
	var ret model.UserKyc
	if err := r.db.Get(&ret, "SELECT * FROM userKyc WHERE userId = ? ORDER BY id DESC LIMIT 1", userId); err != nil {
		return model.UserKyc{}, err
	}
	return ret, nil
}

func (r *UserKycRepository) Get(id int64) (model.UserKyc, error) {
	var ret model.UserKyc
	if err := r.db.Get(&ret, "SELECT * FROM userKyc WHERE id = ?", id); err != nil {
		return model.UserKyc{}, err
	}
	return ret, nil
}

func (r *UserKycRepository) GetAllPending() ([]model.UserKyc, error) {
	var ret []model.UserKyc
	if err := r.db.Select(&ret, "SELECT * FROM userKyc WHERE status = ?", model.UserKycStatusPending); err != nil {
		return nil, err
	}
	return ret, nil
}
