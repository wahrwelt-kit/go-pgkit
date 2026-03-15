package pgutil

import (
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

// IsNoRows reports whether err is or wraps pgx.ErrNoRows.
func IsNoRows(err error) bool {
	return err != nil && errors.Is(err, pgx.ErrNoRows)
}

// IsPgUniqueViolation reports whether err is a PostgreSQL unique violation (code 23505).
func IsPgUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return err != nil && errors.As(err, &pgErr) && pgErr.Code == "23505"
}

// TimestamptzToTime returns a pointer to the time or nil if the timestamptz is invalid.
func TimestamptzToTime(t pgtype.Timestamptz) *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

// TimestamptzToTimeZero returns the time or zero value if the timestamptz is invalid.
func TimestamptzToTimeZero(t pgtype.Timestamptz) time.Time {
	if !t.Valid {
		return time.Time{}
	}
	return t.Time
}

// TimeToTimestamptz converts a time pointer to pgtype.Timestamptz (invalid if t is nil).
func TimeToTimestamptz(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{}
	}
	return pgtype.Timestamptz{Time: *t, Valid: true}
}

// PtrTimeToTime returns the dereferenced time or time.Time{} if t is nil.
func PtrTimeToTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}
