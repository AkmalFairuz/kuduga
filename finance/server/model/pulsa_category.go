package model

type PulsaCategory struct {
	ID         int64  `db:"id"`
	Prefix     string `db:"prefix"`
	CategoryID int64  `db:"categoryId"`
}
