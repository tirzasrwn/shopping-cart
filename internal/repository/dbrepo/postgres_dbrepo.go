package dbrepo

import (
	"database/sql"
	"time"
)

// PostgresDBRepo is the struct used to wrap our database connection pool, so that we
// can easily swap out a real database for a test database, or move to another database
// entirely, as long as the thing being swapped implements all of the functions in the type
// repository.DatabaseRepo.
type PostgresDBRepo struct {
	DB *sql.DB
}

const dbTimeout = time.Second * 3

// Connection returns underlying connection pool.
func (m *PostgresDBRepo) Connection() *sql.DB {
	return m.DB
}
