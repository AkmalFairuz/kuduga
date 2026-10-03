package repository

import (
	"github.com/akmalfairuz/finance/module/database"
	"github.com/akmalfairuz/finance/server/model"
	"time"
)

type TokenRepository struct {
	db *database.DB
}

func NewTokenRepository(db *database.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

func (r *TokenRepository) CreateToken(userId int64, role int, hashedToken []byte, duration time.Duration) (int64, error) {
	insert, err := r.db.Exec(
		"INSERT INTO tokens (userId, role, hashedToken, createdAt, expiredAt, updatedAt) VALUES (?, ?, ?, ?, ?, ?)",
		userId,
		role,
		hashedToken,
		time.Now().Unix(),
		time.Now().Add(duration).Unix(),
		time.Now().Unix(),
	)
	if err != nil {
		return 0, err
	}
	return insert.LastInsertId()
}

func (r *TokenRepository) GetToken(id int64) (model.Token, error) {
	var token model.Token
	if err := r.db.Get(&token, "SELECT * FROM tokens WHERE id = ?", id); err != nil {
		return model.Token{}, err
	}
	return token, nil
}

func (r *TokenRepository) DeleteToken(id int64) error {
	_, err := r.db.Exec("DELETE FROM tokens WHERE id = ?", id)
	return err
}

func (r *TokenRepository) SetDeviceFcmToken(id int64, deviceFcmToken string) error {
	_, err := r.db.Exec("UPDATE tokens SET deviceFcmToken = ?, updatedAt = ? WHERE id = ?", deviceFcmToken, time.Now().Unix(), id)
	return err
}

func (r *TokenRepository) GetTokens(userId int64) ([]model.Token, error) {
	var tokens []model.Token
	if err := r.db.Select(&tokens, "SELECT * FROM tokens WHERE userId = ? ORDER BY updatedAt DESC", userId); err != nil {
		return nil, err
	}
	return tokens, nil
}

func (r *TokenRepository) SetTokenExpiredAt(id int64, expiredAt int64) error {
	if _, err := r.db.Exec("UPDATE tokens SET expiredAt = ?, updatedAt = ? WHERE id = ?", expiredAt, time.Now().Unix(), id); err != nil {
		return err
	}
	return nil
}

func (r *TokenRepository) GetActiveUserCount(from, to int64) (int64, error) {
	var ret struct {
		V int64 `db:"v"`
	}
	if err := r.db.Get(&ret, "SELECT COUNT(DISTINCT userId) AS v FROM tokens WHERE updatedAt >= ? AND updatedAt < ?", from, to); err != nil {
		return 0, err
	}
	return ret.V, nil
}
