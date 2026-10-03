package responsetype

type PaymentMethodResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	NoFee       bool   `json:"noFee"`
	MinAmount   int64  `json:"minAmount"`
	MaxAmount   int64  `json:"maxAmount"`
	ImageUrl    string `json:"imageUrl"`
}
