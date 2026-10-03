package v0

import (
	"fmt"
	"github.com/akmalfairuz/finance/module/database"
)

type V11 struct{}

func (V11) Version() uint {
	return 11
}

func (v V11) Description() string {
	return "Change token to 32 bytes"
}

func (v V11) Up(db *database.DB) error {
	if _, err := db.Exec("DELETE FROM tokens WHERE 1=1"); err != nil {
		return fmt.Errorf("failed to delete all tokens: %w", err)
	}
	if _, err := db.Exec("DELETE FROM otp WHERE 1=1"); err != nil {
		return fmt.Errorf("failed to delete all otp: %w", err)
	}
	if _, err := db.Exec("ALTER TABLE tokens CHANGE COLUMN value hashedToken BINARY(32) NOT NULL"); err != nil {
		return fmt.Errorf("failed to modify tokens.token column: %w", err)
	}
	if _, err := db.Exec("ALTER TABLE otp CHANGE COLUMN token hashedToken BINARY(32) NOT NULL"); err != nil {
		return fmt.Errorf("failed to modify otp.token column: %w", err)
	}
	return nil
}

func (v V11) Down(db *database.DB) error {
	if _, err := db.Exec("DELETE FROM tokens WHERE 1=1"); err != nil {
		return fmt.Errorf("failed to delete all tokens: %w", err)
	}
	if _, err := db.Exec("DELETE FROM otp WHERE 1=1"); err != nil {
		return fmt.Errorf("failed to delete all otp: %w", err)
	}
	if _, err := db.Exec("ALTER TABLE tokens CHANGE COLUMN hashedToken value VARCHAR(255) NOT NULL"); err != nil {
		return fmt.Errorf("failed to modify tokens.token column: %w", err)
	}
	if _, err := db.Exec("ALTER TABLE otp CHANGE COLUMN hashedToken token VARCHAR(255) NOT NULL"); err != nil {
		return fmt.Errorf("failed to modify otp.token column: %w", err)
	}
	return nil
}
