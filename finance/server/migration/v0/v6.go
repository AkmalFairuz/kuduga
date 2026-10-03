package v0

import (
	"fmt"
	"github.com/akmalfairuz/finance/module/database"
)

type V6 struct{}

func (V6) Version() uint {
	return 6
}

func (v V6) Description() string {
	return "Add purchases.productExternalSku and purchases.productProvider column"
}

func (v V6) Up(db *database.DB) error {
	if _, err := db.Exec("ALTER TABLE purchases ADD COLUMN productExternalSku VARCHAR(255) NOT NULL"); err != nil {
		return fmt.Errorf("failed to add products.productExternalSku column: %w", err)
	}
	if _, err := db.Exec("ALTER TABLE purchases ADD COLUMN productProvider VARCHAR(255) NOT NULL"); err != nil {
		return fmt.Errorf("failed to add products.productProvider column: %w", err)
	}
	return nil
}

func (v V6) Down(db *database.DB) error {
	if _, err := db.Exec("ALTER TABLE purchases DROP COLUMN productExternalSku"); err != nil {
		return fmt.Errorf("failed to drop products.productExternalSku column: %w", err)
	}
	if _, err := db.Exec("ALTER TABLE purchases DROP COLUMN productProvider"); err != nil {
		return fmt.Errorf("failed to drop products.productProvider column: %w", err)
	}
	return nil
}
