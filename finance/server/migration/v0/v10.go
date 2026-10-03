package v0

import (
	"fmt"
	"github.com/akmalfairuz/finance/module/database"
)

type V10 struct{}

func (V10) Version() uint {
	return 10
}

func (v V10) Description() string {
	return "Increase token length, update tokens.value and otp.token column type to VARCHAR(255) from CHAR(32)"
}

func (v V10) Up(db *database.DB) error {
	if _, err := db.Exec("ALTER TABLE tokens MODIFY COLUMN value VARCHAR(255) NOT NULL"); err != nil {
		return fmt.Errorf("failed to modify tokens.token column: %w", err)
	}
	if _, err := db.Exec("ALTER TABLE otp MODIFY COLUMN token VARCHAR(255) NOT NULL"); err != nil {
		return fmt.Errorf("failed to modify otp.token column: %w", err)
	}
	return nil
}

func (v V10) Down(db *database.DB) error {
	if _, err := db.Exec("ALTER TABLE tokens MODIFY COLUMN value CHAR(32) NOT NULL"); err != nil {
		return fmt.Errorf("failed to modify tokens.token column: %w", err)
	}
	if _, err := db.Exec("ALTER TABLE otp MODIFY COLUMN token VARCHAR(32) NOT NULL"); err != nil {
		return fmt.Errorf("failed to modify otp.token column: %w", err)
	}
	return nil
}
