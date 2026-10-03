package repository

import (
	"github.com/akmalfairuz/finance/module/database"
	"github.com/akmalfairuz/finance/server/model"
	"github.com/akmalfairuz/finance/server/model/errortype"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"
)

type UserRepository struct {
	db *database.DB
}

func NewUserRepository(db *database.DB) *UserRepository {
	return &UserRepository{db}
}

var ErrBalanceNotEnough = errortype.ErrBalanceNotEnough

func (r *UserRepository) CreateUser(name, displayName, email string, deviceUniqueId *string, passwordHash []byte) (int64, error) {
	insert, err := r.db.Exec("INSERT INTO users (name, displayName, email, passwordHash, deviceUniqueId, createdAt) VALUES (?, ?, ?, ?, ?, ?)", name, displayName, email, passwordHash, deviceUniqueId, time.Now().Unix())
	if err != nil {
		return 0, err
	}
	return insert.LastInsertId()
}

func (r *UserRepository) GetUser(id int64) (model.User, error) {
	var user model.User
	if err := r.db.Get(&user, "SELECT * FROM users WHERE id = ?", id); err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByName(name string) (model.User, error) {
	var user model.User
	if err := r.db.Get(&user, "SELECT * FROM users WHERE name = ?", name); err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (model.User, error) {
	var user model.User
	if err := r.db.Get(&user, "SELECT * FROM users WHERE email = ?", email); err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (r *UserRepository) SetUserPasswordHash(userId int64, passwordHash []byte) error {
	_, err := r.db.Exec("UPDATE users SET passwordHash = ? WHERE id = ?", passwordHash, userId)
	return err
}

func (r *UserRepository) AddUserBalanceWithTx(tx *sqlx.Tx, userId int64, amount int64) error {
	result, err := tx.Exec("UPDATE users SET balance = balance + ? WHERE id = ? AND balance + ? >= 0", amount, userId, amount)
	if rowsAffected, _ := result.RowsAffected(); rowsAffected < 1 {
		return ErrBalanceNotEnough
	}
	return err
}

func (r *UserRepository) UpdateUserEmail(userId int64, email string) error {
	_, err := r.db.Exec("UPDATE users SET email = ? WHERE id = ?", email, userId)
	return err
}

func (r *UserRepository) UpdateUserDisplayName(userId int64, displayName string) error {
	_, err := r.db.Exec("UPDATE users SET displayName = ? WHERE id = ?", displayName, userId)
	return err
}

func (r *UserRepository) UpdateUserPinHash(userId int64, pinHash *[]byte) error {
	_, err := r.db.Exec("UPDATE users SET pinHash = ? WHERE id = ?", pinHash, userId)
	return err
}

func (r *UserRepository) GetUsers(ids []int64) ([]model.User, error) {
	if len(ids) == 0 {
		return []model.User{}, nil
	}
	var ret []model.User
	if err := r.db.Select(&ret, "SELECT * FROM users WHERE id IN (?"+strings.Repeat(",?", len(ids)-1)+")", database.SliceToArgs(ids)...); err != nil {
		return nil, err
	}
	return ret, nil
}

func (r *UserRepository) GetAndLockUser(tx *sqlx.Tx, id int64) (model.User, error) {
	var ret model.User
	if err := tx.Get(&ret, "SELECT * FROM users WHERE id = ? FOR UPDATE", id); err != nil {
		return model.User{}, err
	}
	return ret, nil
}

func (r *UserRepository) UpdateUserKycStatus(userId int64, status bool) error {
	_, err := r.db.Exec("UPDATE users SET kycStatus = ? WHERE id = ?", status, userId)
	return err
}

func (r *UserRepository) GetNewUserCount(from, to int64) (int64, error) {
	var ret struct {
		V int64 `db:"v"`
	}
	if err := r.db.Get(&ret, "SELECT COUNT(*) AS v FROM users WHERE createdAt >= ? AND createdAt < ?", from, to); err != nil {
		return 0, err
	}
	return ret.V, nil
}

func (r *UserRepository) GetTotalUserBalance() (int64, error) {
	var ret struct {
		V int64 `db:"v"`
	}
	if err := r.db.Get(&ret, "SELECT SUM(balance) AS v FROM users"); err != nil {
		return 0, err
	}
	return ret.V, nil
}

func (r *UserRepository) GetUserByReferralCode(referralCode string) (model.User, error) {
	var user model.User
	if err := r.db.Get(&user, "SELECT * FROM users WHERE referralCode = ?", referralCode); err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (r *UserRepository) SetUserReferralTo(userId int64, toUserId int64) error {
	_, err := r.db.Exec("UPDATE users SET referralUserId = ?, referralAt = ? WHERE id = ?", toUserId, time.Now().Unix(), userId)
	return err
}

func (r *UserRepository) SetUserReferralCode(userId int64, referralCode string) error {
	_, err := r.db.Exec("UPDATE users SET referralCode = ? WHERE id = ? AND referralCode IS NULL", referralCode, userId)
	return err
}

func (r *UserRepository) GetCountUsedThisReferralUser(referralUserId int64) (int64, error) {
	var ret struct {
		V int64 `db:"v"`
	}
	if err := r.db.Get(&ret, "SELECT COUNT(*) AS v FROM users WHERE referralUserId = ?", referralUserId); err != nil {
		return 0, err
	}
	return ret.V, nil
}

func (r *UserRepository) UpdateUserDeviceUniqueID(id int64, deviceUniqueID string) error {
	_, err := r.db.Exec("UPDATE users SET deviceUniqueId = ? WHERE id = ?", deviceUniqueID, id)
	return err
}

func (r *UserRepository) GetUserCountByDeviceUniqueIDLast7Days(deviceUniqueID string) (int64, error) {
	var ret struct {
		V int64 `db:"v"`
	}
	if err := r.db.Get(&ret, "SELECT COUNT(*) AS v FROM users WHERE deviceUniqueId = ? AND createdAt >= UNIX_TIMESTAMP()-(86400*7)", deviceUniqueID); err != nil {
		return 0, err
	}
	return ret.V, nil
}

func (r *UserRepository) GetUsersWithBalance() ([]model.User, error) {
	var ret []model.User
	if err := r.db.Select(&ret, "SELECT * FROM users WHERE balance > 0"); err != nil {
		return nil, err
	}
	return ret, nil
}
