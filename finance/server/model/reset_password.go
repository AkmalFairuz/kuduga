package model

type ResetPassword struct {
	ID          int64  `db:"id"`
	UserID      int64  `db:"userId"`
	HashedToken []byte `db:"hashedToken"`
	ResetsAt    *int64 `db:"resetsAt"`
	CreatedAt   int64  `db:"createdAt"`
	ExpiredAt   int64  `db:"expiredAt"`
}
