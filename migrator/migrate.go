package migrator

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrate runs golang-migrate up from migrationsPath (file://) using connURL; ErrNoChange is ignored.
func RunMigrate(connURL, migrationsPath string) error {
	m, err := migrate.New("file://"+migrationsPath, connURL)
	if err != nil {
		return fmt.Errorf("migrator.RunMigrate: New: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migrator.RunMigrate: Up: %w", err)
	}
	return nil
}
