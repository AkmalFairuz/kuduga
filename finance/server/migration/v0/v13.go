package v0

import (
	"fmt"
	"github.com/akmalfairuz/finance/module/database"
)

type V13 struct{}

func (V13) Version() uint {
	return 13
}

func (v V13) Description() string {
	return "Add users.locked"
}

func (v V13) Up(db *database.DB) error {
	if _, err := db.Exec("ALTER TABLE users ADD COLUMN locked BOOLEAN NOT NULL DEFAULT FALSE"); err != nil {
		return fmt.Errorf("failed to add users.locked column: %w", err)
	}
	return nil
}

func (v V13) Down(db *database.DB) error {
	if _, err := db.Exec("ALTER TABLE users DROP COLUMN locked"); err != nil {
		return fmt.Errorf("failed to drop users.locked column: %w", err)
	}
	return nil
}
