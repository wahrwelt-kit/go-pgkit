package migrator

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

// RunFromConn runs goose up migrations from migrationsPath using the given connection string.
func RunFromConn(connStr, migrationsPath string) error {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return fmt.Errorf("migrator.RunFromConn: sql.Open: %w", err)
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("migrator.RunFromConn: SetDialect: %w", err)
	}
	if err := goose.Up(db, migrationsPath); err != nil {
		return fmt.Errorf("migrator.RunFromConn: goose.Up: %w", err)
	}
	return nil
}
