package v0

import (
	"fmt"
	"github.com/akmalfairuz/finance/module/database"
)

type V12 struct{}

func (V12) Version() uint {
	return 12
}

func (v V12) Description() string {
	return "Add token.updatedAt column"
}

func (v V12) Up(db *database.DB) error {
	if _, err := db.Exec("ALTER TABLE tokens ADD COLUMN updatedAt BIGINT NOT NULL"); err != nil {
		return fmt.Errorf("failed to add tokens.updatedAt column: %w", err)
	}
	return nil
}

func (v V12) Down(db *database.DB) error {
	if _, err := db.Exec("ALTER TABLE tokens DROP COLUMN updatedAt"); err != nil {
		return fmt.Errorf("failed to drop tokens.updatedAt column: %w", err)
	}
	return nil
}
