package requesttype

type PaymentCallbackBcaMutasiRequest struct {
	ApiKey  string `json:"apiKey" validate:"required"`
	ID      string `json:"id" validate:"required"`
	Channel string `json:"channel" validate:"required"`
	RefId   string `json:"refId" validate:"required"`
	Status  int    `json:"status" validate:"required"`
	Amount  int64  `json:"amount" validate:"required"`
}
