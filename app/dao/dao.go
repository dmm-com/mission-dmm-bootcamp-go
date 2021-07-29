package dao

import (
	"fmt"
	"log"
	"yatter-backend-go/app/domain/repository"

	"github.com/jmoiron/sqlx"
)

type (
	// DAO interface
	Dao interface {
		// Get account repository
		Account() repository.Account

		// Clear all data in DB
		InitAll() error
	}

	// Implementation for DAO
	dao struct {
		db *sqlx.DB
	}
)

// Create DAO
func New(config DBConfig) (Dao, error) {
	db, err := initDb(config)
	if err != nil {
		return nil, err
	}

	return &dao{db: db}, nil
}

func (d *dao) Account() repository.Account {
	return NewAccount(d.db)
}

func (d *dao) InitAll() error {
	if err := d.exec("SET FOREIGN_KEY_CHECKS=0"); err != nil {
		return fmt.Errorf("Can't disable FOREIGN_KEY_CHECKS: %w", err)
	}

	defer func() {
		err := d.exec("SET FOREIGN_KEY_CHECKS=0")
		if err != nil {
			log.Printf("Can't restore FOREIGN_KEY_CHECKS: %+v", err)
		}
	}()

	for _, table := range []string{"account", "status"} {
		if err := d.exec("TRUNCATE TABLE " + table); err != nil {
			return fmt.Errorf("Can't truncate table "+table+": %w", err)
		}
	}

	return nil
}

func (d *dao) exec(query string, args ...interface{}) error {
	_, err := d.db.Exec(query, args...)
	return err
}
