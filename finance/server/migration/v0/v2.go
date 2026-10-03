package v0

import (
	"fmt"
	"github.com/akmalfairuz/finance/module/database"
)

type V2 struct{}

func (V2) Version() uint {
	return 2
}

func (v V2) Description() string {
	return "Add products.kind and products.purchaseNote & Remove productCategories.dedicated"
}

func (v V2) Up(db *database.DB) error {
	if _, err := db.Exec("ALTER TABLE productCategories DROP COLUMN dedicated"); err != nil {
		return fmt.Errorf("failed to drop productCategories.dedicated column: %w", err)
	}
	if _, err := db.Exec("ALTER TABLE products ADD COLUMN kind VARCHAR(255) NOT NULL, ADD COLUMN purchaseNote TEXT NOT NULL"); err != nil {
		return fmt.Errorf("failed to add product.kind and product.purchaseNote column: %w", err)
	}
	if _, err := db.Exec("ALTER TABLE purchases ADD COLUMN productKind VARCHAR(255) NOT NULL"); err != nil {
		return fmt.Errorf("failed to add purchases.productKind column: %w", err)
	}
	return nil
}

func (v V2) Down(db *database.DB) error {
	if _, err := db.Exec("ALTER TABLE productCategories ADD COLUMN dedicated BOOLEAN NOT NULL"); err != nil {
		return fmt.Errorf("failed to add productCategories.dedicated column: %w", err)
	}
	if _, err := db.Exec("ALTER TABLE products DROP COLUMN kind, DROP COLUMN purchaseNote"); err != nil {
		return fmt.Errorf("failed to drop product.kind and product.purchaseNote column: %w", err)
	}
	if _, err := db.Exec("ALTER TABLE purchases DROP COLUMN productKind"); err != nil {
		return fmt.Errorf("failed to drop purchases.productKind column: %w", err)
	}
	return nil
}
