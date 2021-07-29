package dao

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Interface of configureation
type DBConfig interface {
	FormatDSN() string
}

// Prepare sqlx.DB
func initDb(config DBConfig) (*sqlx.DB, error) {
	driverName := "mysql"
	db, err := sqlx.Open(driverName, config.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("sqlx.Open failed: %w", err)
	}

	return db, nil
}
