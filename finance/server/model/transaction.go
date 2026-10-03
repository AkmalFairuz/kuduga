package model

const (
	_ = iota
	TransactionTypeDeposit
	TransactionTypeWithdraw
	TransactionTypeTransfer
	TransactionTypePurchase
	TransactionTypeTransferToUser
	TransactionTypeReceiveTransferFromUser
	TransactionTypeFee
)

type Transaction struct {
	ID     int64 `db:"id"`
	UserID int64 `db:"userId"`

	Type      int    `db:"type"`
	ExtraData string `db:"extraData"`

	BeforeBalance int64 `db:"beforeBalance"`
	AfterBalance  int64 `db:"afterBalance"`

	Description string `db:"description"`
	Amount      int64  `db:"amount"`
	CreatedAt   int64  `db:"createdAt"`
}

type CreateTransaction struct {
	UserID      int64
	Type        int
	IsAdd       bool
	ExtraData   string
	Description string
	Amount      int64
}
