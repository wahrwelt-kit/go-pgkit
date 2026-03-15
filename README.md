# go-pgkit

PostgreSQL helpers for pgx: error checks, timestamptz converters, pool with retry, and migration runners (goose and golang-migrate).

## Install

```bash
go get github.com/TakuyaYagam1/go-pgkit
```

```go
import "github.com/TakuyaYagam1/go-pgkit/pgutil"
import "github.com/TakuyaYagam1/go-pgkit/postgres"
import "github.com/TakuyaYagam1/go-pgkit/migrator"
```

## Subpackages

### pgutil

- **IsNoRows(err)** — true if err is or wraps pgx.ErrNoRows
- **IsPgUniqueViolation(err)** — true if PostgreSQL unique violation (23505)
- **TimestamptzToTime(t)** — *time.Time or nil if invalid
- **TimestamptzToTimeZero(t)** — time.Time or zero if invalid
- **TimeToTimestamptz(t)** — pgtype.Timestamptz (invalid if t is nil)
- **PtrTimeToTime(t)** — dereference or time.Time{}

### postgres

- **Config** — URL, MaxConns, MinConns
- **New(cfg)** — create pgxpool with exponential backoff retry until Ping succeeds

### migrator

Two runners; use the one that matches your migration layout.

- **RunFromConn(connStr, migrationsPath)** — goose: SQL files with `-- +goose Up` / `-- +goose Down`
- **RunMigrate(connURL, migrationsPath)** — golang-migrate: `NNNNNN_name.up.sql` / `NNNNNN_name.down.sql`; treats ErrNoChange as success

## Example

```go
pool, err := postgres.New(&postgres.Config{
    URL:     os.Getenv("DATABASE_URL"),
    MaxConns: 20,
})
if err != nil {
    log.Fatal(err)
}
defer pool.Close()

if err := migrator.RunFromConn(connStr, "./migrations"); err != nil {
    log.Fatal(err)
}

var t pgtype.Timestamptz
err := pool.QueryRow(ctx, "SELECT created_at FROM users WHERE id = $1", id).Scan(&t)
if pgutil.IsNoRows(err) {
    return ErrNotFound
}
createdAt := pgutil.TimestamptzToTime(t)
```
