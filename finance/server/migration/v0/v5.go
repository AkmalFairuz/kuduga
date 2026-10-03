package v0

import (
	"fmt"
	"github.com/akmalfairuz/finance/module/database"
)

type V5 struct{}

func (V5) Version() uint {
	return 5
}

func (v V5) Description() string {
	return "Add products.proofParser and purchases.parsedProof column"
}

func (v V5) Up(db *database.DB) error {
	if _, err := db.Exec("ALTER TABLE products ADD COLUMN proofParser VARCHAR(255) NOT NULL"); err != nil {
		return fmt.Errorf("failed to add products.proofParser column: %w", err)
	}
	if _, err := db.Exec("ALTER TABLE purchases ADD COLUMN parsedProof VARCHAR(255) NOT NULL"); err != nil {
		return fmt.Errorf("failed to add purchases.parsedProof column: %w", err)
	}
	return nil
}

func (v V5) Down(db *database.DB) error {
	if _, err := db.Exec("ALTER TABLE products DROP COLUMN proofParser"); err != nil {
		return fmt.Errorf("failed to drop products.proofParser column: %w", err)
	}
	if _, err := db.Exec("ALTER TABLE purchases DROP COLUMN parsedProof"); err != nil {
		return fmt.Errorf("failed to drop purchases.parsedProof column: %w", err)
	}
	return nil
}
