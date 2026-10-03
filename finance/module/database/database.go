package database

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"
)

type DB struct {
	*sqlx.DB
}

func New(dsn string, driver string) (*DB, error) {
	attempt := 0
	for {
		db, err := sqlx.Open(driver, dsn)
		if err != nil {
			return nil, err
		}
		if err := db.Ping(); err != nil {
			attempt++
			if attempt > 5 {
				return nil, fmt.Errorf("failed to connect to database after 5 attempts: %w", err)
			}
			time.Sleep(time.Second * 3) // delay 3 seconds
			continue
		}
		return &DB{db}, nil
	}
}

func (db *DB) Exists(query string, args ...any) (bool, error) {
	count, err := db.Count(query, args...)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (db *DB) Count(query string, args ...any) (int, error) {
	rows, err := db.Queryx(query, args...)
	if err != nil {
		return 0, err
	}
	count := 0
	for rows.Next() {
		count++
	}

	if err := rows.Err(); err != nil {
		return 0, err
	}
	return count, nil
}

func (db *DB) DropTable(tableNames ...string) error {
	for _, tableName := range tableNames {
		if _, err := db.Exec("DROP TABLE " + tableName); err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) BulkInsert(table string, columns []string, values ...[]any) (sql.Result, error) {
	return BulkInsert(db, table, columns, values)
}

func BulkInsert(db sqlx.Execer, table string, columns []string, values [][]any) (sql.Result, error) {
	if len(values) == 0 {
		return driver.RowsAffected(0), nil
	}
	rq := " (?" + strings.Repeat(",?", len(columns)-1) + ")"
	args := make([]any, 0, len(columns)*len(values))
	for _, value := range values {
		args = append(args, value...)
	}
	rargs := ""
	for i := 0; i < len(values); i++ {
		rargs += rq
		if i != len(values)-1 {
			rargs += ", "
		}
	}
	res, err := db.Exec("INSERT INTO "+table+" ("+strings.Join(columns, ", ")+") VALUES "+rargs, args...)
	if err != nil {
		return nil, err
	}
	return res, nil
}
