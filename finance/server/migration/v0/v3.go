package v0

import (
	"fmt"
	"github.com/akmalfairuz/finance/module/database"
)

type V3 struct{}

func (V3) Version() uint {
	return 3
}

func (v V3) Description() string {
	return "Add purchases.userSellPrice column"
}

func (v V3) Up(db *database.DB) error {
	if _, err := db.Exec("ALTER TABLE purchases ADD COLUMN userSellPrice BIGINT NOT NULL DEFAULT 0"); err != nil {
		return fmt.Errorf("failed to add purchases.userSellPrice column: %w", err)
	}
	return nil
}

func (v V3) Down(db *database.DB) error {
	if _, err := db.Exec("ALTER TABLE purchases DROP COLUMN userSellPrice"); err != nil {
		return fmt.Errorf("failed to drop purchases.userSellPrice column: %w", err)
	}
	return nil
}
