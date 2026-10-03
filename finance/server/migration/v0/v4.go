package v0

import (
	"fmt"
	"github.com/akmalfairuz/finance/module/database"
)

type V4 struct{}

func (V4) Version() uint {
	return 4
}

func (v V4) Description() string {
	return "Add purchases.note column"
}

func (v V4) Up(db *database.DB) error {
	if _, err := db.Exec("ALTER TABLE purchases ADD COLUMN note TEXT NOT NULL"); err != nil {
		return fmt.Errorf("failed to add purchases.note column: %w", err)
	}
	return nil
}

func (v V4) Down(db *database.DB) error {
	if _, err := db.Exec("ALTER TABLE purchases DROP COLUMN note"); err != nil {
		return fmt.Errorf("failed to drop purchases.note column: %w", err)
	}
	return nil
}
