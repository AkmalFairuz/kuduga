package model

import "time"

const (
	DepositStatusWaiting = iota
	DepositStatusSuccess
	DepositStatusFailed
)

type DepositRequest struct {
	ID                  int64  `db:"id"`
	UserID              int64  `db:"userId"`
	TransactionID       *int64 `db:"transactionId"`
	PaymentMethod       int    `db:"paymentMethod"`
	ExternalPaymentId   string `db:"externalPaymentId"`
	ExternalPaymentData string `db:"externalPaymentData"`
	Amount              int64  `db:"amount"`
	Fee                 int64  `db:"fee"`
	TotalAmount         int64  `db:"totalAmount"`
	Description         string `db:"description"`
	Status              int    `db:"status"`
	CreatedAt           int64  `db:"createdAt"`
	ExpiredAt           int64  `db:"expiredAt"`
}

type DepositRequestTrack struct {
	ID               int64  `db:"id"`
	DepositRequestID int64  `db:"depositRequestId"`
	Description      string `db:"description"`
	NewStatus        *int   `db:"newStatus"`
	CreatedAt        int64  `db:"createdAt"`
}

type CreateDepositRequest struct {
	UserID              int64
	ExternalPaymentId   string
	ExternalPaymentData string
	Amount              int64
	Fee                 int64
	Description         string
	PaymentMethod       int
	ExpireIn            time.Duration
}
