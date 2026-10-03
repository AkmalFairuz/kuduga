package repository

import (
	"github.com/akmalfairuz/finance/module/database"
	"github.com/akmalfairuz/finance/module/pointer"
	"github.com/akmalfairuz/finance/server/model"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"
)

type DepositRepository struct {
	db *database.DB
}

func NewDepositRepository(db *database.DB) *DepositRepository {
	return &DepositRepository{
		db: db,
	}
}

func (r *DepositRepository) CreateDeposit(create model.CreateDepositRequest) (int64, error) {
	insert, err := r.db.Exec("INSERT INTO depositRequests (userId, paymentMethod, externalPaymentId, externalPaymentData, description, amount, fee, totalAmount, status, createdAt, expiredAt) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", create.UserID, create.PaymentMethod, create.ExternalPaymentId, create.ExternalPaymentData, create.Description, create.Amount, create.Fee, create.Amount-create.Fee, model.DepositStatusWaiting, time.Now().Unix(), time.Now().Add(create.ExpireIn).Unix())
	if err != nil {
		return 0, err
	}
	return insert.LastInsertId()
}

func (r *DepositRepository) GetDeposit(id int64) (model.DepositRequest, error) {
	var depositRequest model.DepositRequest
	err := r.db.Get(&depositRequest, "SELECT * FROM depositRequests WHERE id = ?", id)
	if err != nil {
		return model.DepositRequest{}, err
	}
	return depositRequest, nil
}

func (r *DepositRepository) GetDepositByExternalPaymentId(externalPaymentId string, paymentMethod int) (model.DepositRequest, error) {
	var depositRequest model.DepositRequest
	err := r.db.Get(&depositRequest, "SELECT * FROM depositRequests WHERE externalPaymentId = ? AND paymentMethod = ?", externalPaymentId, paymentMethod)
	if err != nil {
		return model.DepositRequest{}, err
	}
	return depositRequest, nil
}

func (r *DepositRepository) SetDepositSuccessWithTx(tx *sqlx.Tx, id int64, description string, txnId int64) error {
	if _, err := tx.Exec("UPDATE depositRequests SET status = ?, transactionId = ? WHERE id = ?", model.DepositStatusSuccess, txnId, id); err != nil {
		return err
	}

	if _, err := tx.Exec("INSERT INTO depositRequestTracks (depositRequestId, description, newStatus, createdAt) VALUES (?, ?, ?, ?)", id, description, model.DepositStatusSuccess, time.Now().Unix()); err != nil {
		return err
	}

	return nil
}

func (r *DepositRepository) SetDepositStatus(id int64, status int, description string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	if _, err := tx.Exec("UPDATE depositRequests SET status = ? WHERE id = ?", status, id); err != nil {
		_ = tx.Rollback()
		return err
	}

	if _, err := tx.Exec("INSERT INTO depositRequestTracks (depositRequestId, description, newStatus, createdAt) VALUES (?, ?, ?, ?)", id, description, status, time.Now().Unix()); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *DepositRepository) GetDepositTracks(depositRequestId int64) ([]model.DepositRequestTrack, error) {
	var depositRequestTracks []model.DepositRequestTrack
	err := r.db.Select(&depositRequestTracks, "SELECT * FROM depositRequestTracks WHERE depositRequestId = ? ORDER BY createdAt DESC", depositRequestId)
	if err != nil {
		return nil, err
	}
	return depositRequestTracks, nil
}

func (r *DepositRepository) GetDepositByUserId(userId int64) ([]model.DepositRequest, error) {
	var depositRequests []model.DepositRequest
	err := r.db.Select(&depositRequests, "SELECT * FROM depositRequests WHERE userId = ? ORDER BY createdAt DESC", userId)
	if err != nil {
		return nil, err
	}
	return depositRequests, nil
}

func (r *DepositRepository) GetDepositByUserIdAndStatus(userId int64, status int) ([]model.DepositRequest, error) {
	var depositRequests []model.DepositRequest
	err := r.db.Select(&depositRequests, "SELECT * FROM depositRequests WHERE userId = ? AND status = ? ORDER BY createdAt DESC", userId, status)
	if err != nil {
		return nil, err
	}
	return depositRequests, nil
}

func (r *DepositRepository) GetDepositIdByExternalPaymentId(paymentMethodId int, externalPaymentId string) (int64, error) {
	var id struct {
		ID int64 `db:"id"`
	}
	if err := r.db.Get(&id, "SELECT id FROM depositRequests WHERE externalPaymentId = ? AND paymentMethod = ?", externalPaymentId, paymentMethodId); err != nil {
		return 0, err
	}
	return id.ID, nil
}

func (r *DepositRepository) Begin() (*sqlx.Tx, error) {
	return r.db.Beginx()
}

func (r *DepositRepository) GetDepositWillExpire() ([]model.DepositRequest, error) {
	var ret []model.DepositRequest
	if err := r.db.Select(&ret, "SELECT * FROM depositRequests WHERE status = ? AND ? >= expiredAt", model.DepositStatusWaiting, time.Now().Unix()); err != nil {
		return nil, err
	}
	return ret, nil
}

func (r *DepositRepository) GetAllDepositCount(from, to int64) (int64, error) {
	var ret struct {
		V int64 `db:"v"`
	}
	if err := r.db.Get(&ret, "SELECT COUNT(*) AS v FROM depositRequests WHERE status = ? AND createdAt >= ? AND createdAt < ?", model.DepositStatusSuccess, from, to); err != nil {
		return 0, err
	}
	return ret.V, nil
}

func (r *DepositRepository) GetTotalDepositAmount(from, to int64) (int64, error) {
	var ret struct {
		V int64 `db:"v"`
	}
	if err := r.db.Get(&ret, "SELECT COALESCE(SUM(amount) - SUM(fee), 0) AS v FROM depositRequests WHERE status = ? AND createdAt >= ? AND createdAt < ?", model.DepositStatusSuccess, from, to); err != nil {
		return 0, err
	}
	return ret.V, nil
}

func (r *DepositRepository) SetDepositExpire(deposits []model.DepositRequest, description string) error {
	if len(deposits) == 0 {
		return nil
	}
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	args := []int64{model.DepositStatusFailed}
	for _, deposit := range deposits {
		args = append(args, deposit.ID)
	}
	args = append(args, model.DepositStatusWaiting)
	return database.Tx(tx, func() error {
		if _, err := tx.Exec("UPDATE depositRequests SET status = ? WHERE id IN (?"+strings.Repeat(",?", len(deposits)-1)+") AND status = ?", database.SliceToArgs(args)...); err != nil {
			return err
		}

		args2 := make([][]any, 0)
		for _, deposit := range deposits {
			args2 = append(args2, []any{deposit.ID, description, pointer.Make(model.DepositStatusFailed), time.Now().Unix()})
		}

		if _, err := database.BulkInsert(tx, "depositRequestTracks", []string{"depositRequestId", "description", "newStatus", "createdAt"}, args2); err != nil {
			return err
		}
		return nil
	})
}
