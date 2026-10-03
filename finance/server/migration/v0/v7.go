package v0

import (
	"fmt"
	"github.com/akmalfairuz/finance/module/database"
)

type V7 struct{}

func (V7) Version() uint {
	return 7
}

func (v V7) Description() string {
	return "Add notifications.refId column"
}

func (v V7) Up(db *database.DB) error {
	if _, err := db.Exec("ALTER TABLE notifications ADD COLUMN refId VARCHAR(255) NOT NULL"); err != nil {
		return fmt.Errorf("failed to add notifications.refId column: %w", err)
	}
	return nil
}

func (v V7) Down(db *database.DB) error {
	if _, err := db.Exec("ALTER TABLE notifications DROP COLUMN refId"); err != nil {
		return fmt.Errorf("failed to drop notifications.refId column: %w", err)
	}
	return nil
}
