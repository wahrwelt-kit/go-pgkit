package migrator

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func startPostgres(t *testing.T) (connStr string, cleanup func()) {
	t.Helper()
	ctx := context.Background()
	c, err := postgres.Run(ctx, "postgres:17-alpine",
		postgres.WithDatabase("pgkit_test"),
		postgres.WithUsername("u"),
		postgres.WithPassword("p"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2)),
	)
	require.NoError(t, err)
	connStr, err = c.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)
	return connStr, func() { _ = c.Terminate(ctx) }
}

func TestRunFromConn(t *testing.T) {
	connStr, cleanup := startPostgres(t)
	defer cleanup()

	goosePath, err := filepath.Abs("testdata/goose")
	require.NoError(t, err)
	require.NoError(t, RunFromConn(connStr, goosePath))

	pool, err := pgxpool.New(context.Background(), connStr)
	require.NoError(t, err)
	defer pool.Close()
	var n int
	err = pool.QueryRow(context.Background(), "SELECT 1 FROM pg_tables WHERE tablename = 'pgkit_test'").Scan(&n)
	require.NoError(t, err)
	require.Equal(t, 1, n)
}

func TestRunMigrate(t *testing.T) {
	connStr, cleanup := startPostgres(t)
	defer cleanup()

	migratePath, err := filepath.Abs("testdata/migrate")
	require.NoError(t, err)
	require.NoError(t, RunMigrate(connStr, migratePath))

	pool, err := pgxpool.New(context.Background(), connStr)
	require.NoError(t, err)
	defer pool.Close()
	var n int
	err = pool.QueryRow(context.Background(), "SELECT 1 FROM pg_tables WHERE tablename = 'pgkit_migrate_test'").Scan(&n)
	require.NoError(t, err)
	require.Equal(t, 1, n)
}
