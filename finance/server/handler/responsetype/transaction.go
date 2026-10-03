package responsetype

type TransactionInfo struct {
	Id            int64  `json:"id"`
	Amount        int64  `json:"amount"`
	Description   string `json:"description"`
	BeforeBalance int64  `json:"beforeBalance"`
	AfterBalance  int64  `json:"afterBalance"`
	CreatedAt     int64  `json:"createdAt"`
}
