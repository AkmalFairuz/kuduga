package v0

import (
	"github.com/akmalfairuz/finance/module/database"
)

type V17 struct{}

func (V17) Version() uint {
	return 17
}

func (v V17) Description() string {
	return "Add user.deviceUniqueId"
}

func (v V17) Up(db *database.DB) error {
	if _, err := db.Exec(`ALTER TABLE users ADD COLUMN deviceUniqueId VARCHAR(255)`); err != nil {
		return err
	}
	return nil
}

func (v V17) Down(db *database.DB) error {
	if _, err := db.Exec("ALTER TABLE users DROP COLUMN deviceUniqueId"); err != nil {
		return err
	}
	return nil
}
