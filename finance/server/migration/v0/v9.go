package v0

import (
	"github.com/akmalfairuz/finance/module/database"
)

type V9 struct{}

func (V9) Version() uint {
	return 9
}

func (v V9) Description() string {
	return "Add products.billAdmin, purchases.billAdmin column, add billPrePurchases.billAdmin column"
}

func (v V9) Up(db *database.DB) error {
	if _, err := db.Exec("ALTER TABLE products ADD COLUMN billAdmin BIGINT NOT NULL DEFAULT 0"); err != nil {
		return err
	}
	if _, err := db.Exec("ALTER TABLE purchases ADD COLUMN billAdmin BIGINT NOT NULL DEFAULT 0"); err != nil {
		return err
	}
	if _, err := db.Exec("ALTER TABLE billPrePurchases ADD COLUMN billAdmin BIGINT NOT NULL DEFAULT 0"); err != nil {
		return err
	}
	return nil
}

func (v V9) Down(db *database.DB) error {
	if _, err := db.Exec("ALTER TABLE products DROP COLUMN billAdmin"); err != nil {
		return err
	}
	if _, err := db.Exec("ALTER TABLE purchases DROP COLUMN billAdmin"); err != nil {
		return err
	}
	if _, err := db.Exec("ALTER TABLE billPrePurchases DROP COLUMN billAdmin"); err != nil {
		return err
	}
	return nil
}
