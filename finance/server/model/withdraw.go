package model

const (
	WithdrawStatusWaiting = iota
	WithdrawStatusSuccess
	WithdrawStatusFailed
)

type WithdrawRequest struct {
	ID                int64  `db:"id"`
	UserID            int64  `db:"userId"`
	TransactionID     *int64 `db:"transactionId"`
	Amount            int64  `db:"amount"`
	Fee               int64  `db:"fee"`
	Status            int    `db:"status"`
	BankType          int    `db:"bankType"`
	BankAccountNumber string `db:"bankAccountNumber"`
	CreatedAt         int64  `db:"createdAt"`
}

type WithdrawRequestTrack struct {
	ID                int64  `db:"id"`
	WithdrawRequestID int64  `db:"withdrawRequestId"`
	Description       string `db:"description"`
	NewStatus         *int   `db:"newStatus"`
	CreatedAt         int64  `db:"createdAt"`
}
