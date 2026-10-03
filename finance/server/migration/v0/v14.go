package v0

import (
	"fmt"
	"github.com/akmalfairuz/finance/module/database"
)

type V14 struct{}

func (V14) Version() uint {
	return 14
}

func (v V14) Description() string {
	return "Add resetPasswords.hashedToken"
}

func (v V14) Up(db *database.DB) error {
	if _, err := db.Exec("DELETE FROM resetPasswords WHERE 1=1"); err != nil {
		return fmt.Errorf("failed to delete resetPasswords: %w", err)
	}
	if _, err := db.Exec("ALTER TABLE resetPasswords DROP COLUMN token"); err != nil {
		return fmt.Errorf("failed to drop resetPasswords.token column: %w", err)
	}
	if _, err := db.Exec("ALTER TABLE resetPasswords ADD COLUMN hashedToken BINARY(32) NOT NULL"); err != nil {
		return fmt.Errorf("failed to add resetPasswords.hashedToken column: %w", err)
	}
	return nil
}

func (v V14) Down(db *database.DB) error {
	if _, err := db.Exec("DELETE FROM resetPasswords WHERE 1=1"); err != nil {
		return fmt.Errorf("failed to delete resetPasswords: %w", err)
	}
	if _, err := db.Exec("ALTER TABLE resetPasswords ADD COLUMN token CHAR(32) NOT NULL"); err != nil {
		return fmt.Errorf("failed to add resetPasswords.token column: %w", err)
	}
	if _, err := db.Exec("ALTER TABLE resetPasswords DROP COLUMN hashedToken"); err != nil {
		return fmt.Errorf("failed to drop resetPasswords.hashedToken column: %w", err)
	}
	return nil
}
