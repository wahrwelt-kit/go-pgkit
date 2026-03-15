// Package migrator provides migration runners: RunFromConn for goose (SQL files
// with +goose Up/Down) and RunMigrate for golang-migrate (up/down .sql files).
// Use the one that matches your migration layout.
package migrator
