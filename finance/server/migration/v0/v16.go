package v0

import (
	"github.com/akmalfairuz/finance/module/database"
)

type V16 struct{}

func (V16) Version() uint {
	return 16
}

func (v V16) Description() string {
	return "Add user referral columns"
}

func (v V16) Up(db *database.DB) error {
	if _, err := db.Exec(`ALTER TABLE users ADD COLUMN referralCode VARCHAR(255) UNIQUE`); err != nil {
		return err
	}
	if _, err := db.Exec(`ALTER TABLE users ADD COLUMN referralUserId BIGINT`); err != nil {
		return err
	}
	if _, err := db.Exec(`ALTER TABLE users ADD COLUMN referralAt BIGINT`); err != nil {
		return err
	}
	return nil
}

func (v V16) Down(db *database.DB) error {
	if _, err := db.Exec("ALTER TABLE users DROP COLUMN referralCode"); err != nil {
		return err
	}
	if _, err := db.Exec("ALTER TABLE users DROP COLUMN referralUserId"); err != nil {
		return err
	}
	if _, err := db.Exec("ALTER TABLE users DROP COLUMN referralAt"); err != nil {
		return err
	}
	return nil
}
