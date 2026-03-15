package pgutil

import (
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsNoRows(t *testing.T) {
	t.Parallel()
	assert.True(t, IsNoRows(pgx.ErrNoRows))
	assert.True(t, IsNoRows(errors.Join(pgx.ErrNoRows, errors.New("other"))))
	assert.False(t, IsNoRows(nil))
	assert.False(t, IsNoRows(errors.New("other")))
}

func TestIsPgUniqueViolation(t *testing.T) {
	t.Parallel()
	assert.True(t, IsPgUniqueViolation(&pgconn.PgError{Code: "23505"}))
	assert.False(t, IsPgUniqueViolation(&pgconn.PgError{Code: "23503"}))
	assert.False(t, IsPgUniqueViolation(nil))
	assert.False(t, IsPgUniqueViolation(errors.New("other")))
}

func TestTimestamptzToTime(t *testing.T) {
	t.Parallel()
	tt := time.Date(2025, 3, 15, 12, 0, 0, 0, time.UTC)
	out := TimestamptzToTime(pgtype.Timestamptz{Time: tt, Valid: true})
	require.NotNil(t, out)
	assert.True(t, out.Equal(tt))
	assert.Nil(t, TimestamptzToTime(pgtype.Timestamptz{Valid: false}))
}

func TestTimestamptzToTimeZero(t *testing.T) {
	t.Parallel()
	tt := time.Date(2025, 3, 15, 12, 0, 0, 0, time.UTC)
	assert.True(t, TimestamptzToTimeZero(pgtype.Timestamptz{Time: tt, Valid: true}).Equal(tt))
	assert.True(t, TimestamptzToTimeZero(pgtype.Timestamptz{Valid: false}).IsZero())
}

func TestTimeToTimestamptz(t *testing.T) {
	t.Parallel()
	tt := time.Date(2025, 3, 15, 12, 0, 0, 0, time.UTC)
	out := TimeToTimestamptz(&tt)
	assert.True(t, out.Valid)
	assert.True(t, out.Time.Equal(tt))
	out = TimeToTimestamptz(nil)
	assert.False(t, out.Valid)
}

func TestPtrTimeToTime(t *testing.T) {
	t.Parallel()
	tt := time.Date(2025, 3, 15, 12, 0, 0, 0, time.UTC)
	assert.True(t, PtrTimeToTime(&tt).Equal(tt))
	assert.True(t, PtrTimeToTime(nil).IsZero())
}
