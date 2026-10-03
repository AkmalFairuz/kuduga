package v0

import (
	"github.com/akmalfairuz/finance/module/database"
)

type V8 struct{}

func (V8) Version() uint {
	return 8
}

func (v V8) Description() string {
	return "Add products.skuPriority column, add products.sku not unique anymore, add productCategories.meta column and add purchases.successAt column"
}

func (v V8) Up(db *database.DB) error {
	if _, err := db.Exec("ALTER TABLE products ADD COLUMN skuPriority INT NOT NULL DEFAULT 0"); err != nil {
		return err
	}
	if _, err := db.Exec("ALTER TABLE products DROP INDEX sku"); err != nil {
		return err
	}
	if _, err := db.Exec("ALTER TABLE productCategories ADD COLUMN meta MEDIUMTEXT NOT NULL"); err != nil {
		return err
	}
	if _, err := db.Exec("ALTER TABLE purchases ADD COLUMN successAt BIGINT NOT NULL DEFAULT 0"); err != nil {
		return err
	}
	return nil
}

func (v V8) Down(db *database.DB) error {
	if _, err := db.Exec("ALTER TABLE products DROP COLUMN skuPriority"); err != nil {
		return err
	}
	if _, err := db.Exec("ALTER TABLE products ADD UNIQUE INDEX sku (sku)"); err != nil {
		return err
	}
	if _, err := db.Exec("ALTER TABLE productCategories DROP COLUMN meta"); err != nil {
		return err
	}
	if _, err := db.Exec("ALTER TABLE purchases DROP COLUMN successAt"); err != nil {
		return err
	}
	return nil
}
