package responsetype

type DepositDetailedInfo struct {
	Id            int64          `json:"id"`
	Amount        int64          `json:"amount"`
	Fee           int64          `json:"fee"`
	PaymentMethod string         `json:"paymentMethod"`
	PaymentData   any            `json:"paymentData"`
	PaymentUrl    string         `json:"paymentUrl,omitempty"`
	Description   string         `json:"description"`
	Status        int            `json:"status"`
	CreatedAt     int64          `json:"createdAt"`
	ExpiredAt     int64          `json:"expiredAt"`
	Tracks        []DepositTrack `json:"tracks"`
}

type DepositInfo struct {
	Id            int64  `json:"id"`
	Amount        int64  `json:"amount"`
	Fee           int64  `json:"fee"`
	PaymentMethod string `json:"paymentMethod"`
	Status        int    `json:"status"`
	CreatedAt     int64  `json:"createdAt"`
}

type DepositTrack struct {
	Description string `json:"description"`
	NewStatus   *int   `json:"newStatus"`
	CreatedAt   int64  `json:"createdAt"`
}
