package model

type Voucher struct {
	ID        int64  `db:"id"`
	UserID    int64  `db:"userId"`
	Type      string `db:"type"`
	ExpiredAt int64  `db:"expiredAt"`
	CreatedAt int64  `db:"createdAt"`
}
