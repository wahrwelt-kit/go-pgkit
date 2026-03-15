package postgres

// Config holds connection URL and optional pool limits for New.
type Config struct {
	URL      string
	MaxConns int
	MinConns int
}
