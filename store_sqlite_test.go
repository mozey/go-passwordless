package passwordless

import (
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func getConnection() (db *sqlx.DB, err error) {
	db, err = sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		return db, errors.WithStack(err)
	}
	_, err = db.Exec(`create table session (
	uid integer primary key autoincrement,
	token varchar(255) not null unique,
	expires datetime not null,
	created datetime not null default current_timestamp
);`)
	if err != nil {
		return db, errors.WithStack(err)
	}
	return db, nil
}

func TestSQLiteStore(t *testing.T) {
	db, err := getConnection()
	require.NoError(t, err)
	s, err := NewSQLiteStore(db, "", "", "", "")
	require.NoError(t, err)
	require.NotNil(t, s)

	b, exp, err := s.Exists(nil, "uid")
	require.NoError(t, err)
	require.False(t, b)
	require.True(t, exp.IsZero())

	err = s.Store(nil, "", "uid", -time.Hour)
	require.NoError(t, err)
	b, exp, err = s.Exists(nil, "uid")
	require.NoError(t, err)
	require.False(t, b)
	require.True(t, exp.IsZero())

	err = s.Store(nil, "", "uid", time.Hour)
	require.NoError(t, err)
	b, exp, err = s.Exists(nil, "uid")
	require.NoError(t, err)
	require.True(t, b)
	require.False(t, exp.IsZero())
}

func TestSQLiteStoreVerify(t *testing.T) {
	db, err := getConnection()
	require.NoError(t, err)
	s, err := NewSQLiteStore(db, "", "", "", "")
	require.NoError(t, err)
	require.NotNil(t, s)

	// Token doesn't exist
	b, err := s.Verify(nil, "bad_token", "uid")
	require.False(t, b)
	require.Equal(t, ErrTokenNotFound.Error(), err.Error())

	// Token expired
	err = s.Store(nil, "", "uid", -time.Hour)
	require.NoError(t, err)
	b, err = s.Verify(nil, "bad_token", "uid")
	require.False(t, b)
	require.Equal(t, ErrTokenNotFound.Error(), err.Error())

	// Token wrong
	err = s.Store(nil, "token", "uid", time.Hour)
	require.NoError(t, err)
	b, err = s.Verify(nil, "bad_token", "uid")
	require.False(t, b)

	// Token correct
	b, err = s.Verify(nil, "token", "uid")
	require.True(t, b)
	require.NoError(t, err)
}
