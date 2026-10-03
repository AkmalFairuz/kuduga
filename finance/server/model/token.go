package model

import "fmt"

type Token struct {
	ID             int64   `db:"id"`
	UserID         int64   `db:"userId"`
	Role           int     `db:"role"`
	HashedToken    []byte  `db:"hashedToken"`
	DeviceFcmToken *string `db:"deviceFcmToken"`
	CreatedAt      int64   `db:"createdAt"`
	ExpiredAt      int64   `db:"expiredAt"`
	UpdatedAt      int64   `db:"updatedAt"`
}

func (t Token) Format(rawToken string) string {
	return fmt.Sprintf("%d-%s", t.ID, rawToken)
}
