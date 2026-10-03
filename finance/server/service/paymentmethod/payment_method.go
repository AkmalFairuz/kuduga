package paymentmethod

import (
	"time"
)

const (
	BankBcaID = iota + 1
	QrisID
	BrivaID
	BniVirtualID
	AtmBersamaID
	PermataBankVirtualID
	AlfamartID
)

const (
	StatusWaitingPayment = iota + 1
	StatusSuccess
	StatusFailed
)

type PaymentMethod interface {
	NoFee() bool
	ID() int
	Name() string
	ExpireDuration() time.Duration
	CreatePayment(opt CreatePaymentOptions) (Data, error)
	Check(refId string) (CheckResult, error)
	PublicInfo(externalPaymentId, externalPaymentData string) (any, error)
	MinAmount() int64
	MaxAmount() int64
	Description() string
	ImageUrl() string
	Cancel(refId string) error
}

type CreatePaymentOptions struct {
	RefID          string
	CustomerName   string
	Email          string
	ProductDetails string
	Amount         int64
}

type Data struct {
	ExternalPaymentId   string
	ExternalPaymentData string
	Amount              int64
	Fee                 int64
	ExpiredAt           int64
}

type CheckResult struct {
	Status int
}
