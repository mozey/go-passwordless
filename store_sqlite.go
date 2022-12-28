package passwordless

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// SQLiteStore is a Store that keeps tokens in SQLite
type SQLiteStore struct {
	db *sqlx.DB
}

// NewSQLiteStore creates and returns a new SQLiteStore
func NewSQLiteStore(db *sqlx.DB) *SQLiteStore {
	return &SQLiteStore{
		db: db,
	}
}

// Store a generated token in SQLite for a user
func (s SQLiteStore) Store(ctx context.Context, token, uid string, ttl time.Duration) error {
	return errors.Errorf("TODO")
}

// Exists checks to see if a token exists
func (s SQLiteStore) Exists(ctx context.Context, uid string) (bool, time.Time, error) {
	return errors.Errorf("TODO")
}

// Verify checks to see if a token exists and is valid for a user
func (s SQLiteStore) Verify(ctx context.Context, token, uid string) (bool, error) {
	return errors.Errorf("TODO")
}

// Delete removes a key from the store
func (s SQLiteStore) Delete(ctx context.Context, uid string) error {
	return errors.Errorf("TODO")
}
