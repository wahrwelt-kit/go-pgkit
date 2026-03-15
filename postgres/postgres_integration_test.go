package postgres

import (
	"context"
	"testing"

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

func TestNew(t *testing.T) {
	connStr, cleanup := startPostgres(t)
	defer cleanup()

	pool, err := New(&Config{URL: connStr})
	require.NoError(t, err)
	defer pool.Close()
	require.NoError(t, pool.Ping(context.Background()))
}

func TestNew_InvalidURL(t *testing.T) {
	_, err := New(&Config{URL: "postgres://invalid:5432/nonexistent?sslmode=disable"})
	require.Error(t, err)
}
