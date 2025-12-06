// internal/db/db.go
package db

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

var schemaSQL string

// var sqldbsync
var (
	instance *sql.DB
	once     sync.Once
)

// Init initializes the DB once. Call this early (main) with your DSN.
func Init(dsn string) error {
	var err error
	once.Do(func() {
		instance, err = sql.Open("sqlite3", dsn)
		if err == nil {
			err = instance.Ping()
		}
	})
	return err
}

// Connect initializes the package-level DB. Call this once from main.
func Connect(dsn string) error {
	var err error
	once.Do(func() {
		instance, err = sql.Open("sqlite3", dsn)
		if err != nil {
			return
		}
		err = instance.Ping()
	})
	return err
}

// Close closes the package-level DB.
func Close() error {
	if instance == nil {
		return nil
	}
	return instance.Close()
}

// GetDB returns the initialized *sql.DB (may be nil if Init wasn't called).
func GetDB() *sql.DB {
	return instance
}

func InitDB(dbPath string) error {
	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("error al abrir la base de datos: %w", err)
	}
	DB.SetMaxOpenConns(1)

	for _, stmt := range strings.Split(schemaSQL, ";") {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if _, err := DB.Exec(stmt); err != nil {
			return fmt.Errorf("error en la migración: %w", err)
		}
	}
	return nil
}
