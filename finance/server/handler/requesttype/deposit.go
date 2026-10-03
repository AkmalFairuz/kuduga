package requesttype

type DepositCreateRequest struct {
	PaymentMethod int   `form:"paymentMethod"`
	Amount        int64 `form:"amount"`
}

type DepositHistoryRequest struct {
	FilterStatus *int `query:"filterStatus"`
}

type DepositDetailedRequest struct {
	ID int64 `query:"id" validate:"required"`
}
