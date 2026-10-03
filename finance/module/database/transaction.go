package database

import "github.com/jmoiron/sqlx"

func Tx(tx *sqlx.Tx, onTx func() error) error {
	if err := onTx(); err != nil {
		if err2 := tx.Rollback(); err2 != nil {
			return err2
		}
		return err
	}
	return tx.Commit()
}
