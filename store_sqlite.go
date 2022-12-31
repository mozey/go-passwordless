package passwordless

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// SQLiteStore is a Store that keeps tokens in SQLite
type SQLiteStore struct {
	db *sqlx.DB
	// table name for sessions
	table string
	// colToken name for token column
	colToken string
	// colUID name for uid column
	colUID string
	// colExpires name for expires column
	colExpires string
	// dateFormat for colExpires timestamp
	dateFormat string
}

// NewSQLiteStore creates and returns a new SQLiteStore
func NewSQLiteStore(db *sqlx.DB, table, colToken, colUID, colExpires string) (store *SQLiteStore, err error) {
	if db == nil {
		return store, errors.Errorf("invalid db connection")
	}
	if table == "" {
		table = "session"
	}
	if colToken == "" {
		colToken = "token"
	}
	if colUID == "" {
		colUID = "uid"
	}
	if colExpires == "" {
		colExpires = "expires"
	}
	dateFormatMicro := "2006-01-02 15:04:05.000000"
	return &SQLiteStore{
		db:         db,
		table:      table,
		colToken:   colToken,
		colUID:     colUID,
		colExpires: colExpires,
		dateFormat: dateFormatMicro,
	}, nil
}

// Store a generated token in SQLite for a user
func (s SQLiteStore) Store(ctx context.Context, token, uid string, ttl time.Duration) (err error) {
	query := fmt.Sprintf("insert into %s (%s, %s, %s) values (:values)",
		s.table, s.colUID, s.colToken, s.colExpires)

	values := make([]interface{}, 0, 1)
	row := make([]interface{}, 3)
	row[0] = uid
	row[1] = token
	row[2] = time.Now().Add(ttl).Format(s.dateFormat)
	values = append(values, row)

	query, _, err = sqlx.Named(query, map[string]interface{}{
		"values": values,
	})
	if err != nil {
		return errors.WithStack(err)
	}

	query, args, err := sqlx.In(query, values...)
	if err != nil {
		return errors.WithStack(err)
	}

	fmt.Println("query", query)
	fmt.Println(fmt.Sprintf("args %#v", args))
	_, err = s.db.Exec(query, args...)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Exists checks to see if a token exists
func (s SQLiteStore) Exists(ctx context.Context, uid string) (bool, time.Time, error) {
	return false, time.Now(), errors.Errorf("TODO Exists")
}

// Verify checks to see if a token exists and is valid for a user
func (s SQLiteStore) Verify(ctx context.Context, token, uid string) (valid bool, err error) {
	var row map[string]interface{}
	err = namedSelect(s.db, &row,
		"select * from session where uid = :uid and token = :token",
		map[string]interface{}{
			"uid":   uid,
			"token": token,
		})
	if err != nil {
		fmt.Println(err)
	} else {
		b, _ := json.MarshalIndent(row, "", "  ")
		fmt.Println(string(b))
	}
	return false, errors.WithStack(ErrTokenNotFound)
}

// Delete removes a key from the store
func (s SQLiteStore) Delete(ctx context.Context, uid string) error {
	return errors.Errorf("TODO Delete")
}

// namedGet is a wrapper around db.Get for queries using named params
func namedGet(db *sqlx.DB, dest interface{}, query string, arg interface{}) (err error) {
	st, err := db.PrepareNamed(query)
	if err != nil {
		return errors.WithStack(err)
	}
	err = st.Get(dest, arg)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// namedSelect is a wrapper around db.Select for queries using named params
func namedSelect(db *sqlx.DB, dest interface{}, query string, arg interface{}) (err error) {
	st, err := db.PrepareNamed(query)
	if err != nil {
		return errors.WithStack(err)
	}
	err = st.Select(dest, arg)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
