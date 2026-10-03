package repository

import (
	"github.com/akmalfairuz/finance/module/database"
	"github.com/akmalfairuz/finance/module/hash"
	"github.com/akmalfairuz/finance/server/model"
	"time"
)

type OTPRepository struct {
	db *database.DB
}

func NewOTPRepository(db *database.DB) *OTPRepository {
	return &OTPRepository{db: db}
}

func (r *OTPRepository) CreateOTP(otpType, scope int, contact, code, token string, duration time.Duration) (int64, error) {
	//var cnt struct {
	//	Count int `db:"c"`
	//}
	//if err := r.db.Get(&cnt, "SELECT COUNT(*) as c FROM otp WHERE contact = ? AND type = ? AND scope = ?", contact, otpType, scope); err != nil {
	//	return 0, err
	//}
	//if cnt.Count >= 5 {
	//	return 0, errortype.ErrTooManyOTPRequest
	//}

	hashedToken := hash.Sha256Bytes([]byte(token))
	insert, err := r.db.Exec(
		"INSERT INTO otp (type, scope, contact, code, hashedToken, createdAt, expiredAt) VALUES (?, ?, ?, ?, ?, ?, ?)",
		otpType,
		scope,
		contact,
		code,
		hashedToken,
		time.Now().Unix(),
		time.Now().Add(duration).Unix())
	if err != nil {
		return 0, err
	}
	return insert.LastInsertId()
}

func (r *OTPRepository) UseOTP(otp model.OTP) error {
	if _, err := r.db.Exec("UPDATE otp SET usedAt = ? WHERE id = ?", time.Now().Unix(), otp.ID); err != nil {
		return err
	}
	return nil
}

func (r *OTPRepository) Retry(otp model.OTP, newCode string) error {
	if _, err := r.db.Exec("UPDATE otp SET retryAttempt = retryAttempt + 1, lastRetryAt = ?, code = ? WHERE id = ?", time.Now().Unix(), newCode, otp.ID); err != nil {
		return err
	}
	return nil
}

func (r *OTPRepository) GetOTP(id int64) (model.OTP, error) {
	var otp model.OTP
	if err := r.db.Get(&otp, "SELECT * FROM otp WHERE id = ?", id); err != nil {
		return model.OTP{}, err
	}
	return otp, nil
}

func (r *OTPRepository) AddOTPFailAttempt(id int64) error {
	if _, err := r.db.Exec("UPDATE otp SET failAttempt = failAttempt + 1 WHERE id = ?", id); err != nil {
		return err
	}
	return nil
}
