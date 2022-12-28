package passwordless

import (
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSQLiteStore(t *testing.T) {
	var db *sqlx.DB
	ms := NewSQLiteStore(db)
	assert.NotNil(t, ms)

	b, exp, err := ms.Exists(nil, "uid")
	require.NoError(t, err)
	assert.False(t, b)
	assert.True(t, exp.IsZero())

	err = ms.Store(nil, "", "uid", -time.Hour)
	b, exp, err = ms.Exists(nil, "uid")
	require.NoError(t, err)
	assert.False(t, b)
	assert.True(t, exp.IsZero())

	err = ms.Store(nil, "", "uid", time.Hour)
	b, exp, err = ms.Exists(nil, "uid")
	require.NoError(t, err)
	assert.True(t, b)
	assert.False(t, exp.IsZero())
}

func TestSQLiteStoreVerify(t *testing.T) {
	var db *sqlx.DB
	ms := NewSQLiteStore(db)
	assert.NotNil(t, ms)

	// Token doesn't exist
	b, err := ms.Verify(nil, "badtoken", "uid")
	assert.False(t, b)
	assert.Equal(t, ErrTokenNotFound, err)

	// Token expired
	err = ms.Store(nil, "", "uid", -time.Hour)
	b, err = ms.Verify(nil, "badtoken", "uid")
	assert.False(t, b)
	assert.Equal(t, ErrTokenNotFound, err)

	// Token wrong
	err = ms.Store(nil, "token", "uid", time.Hour)
	b, err = ms.Verify(nil, "badtoken", "uid")
	assert.False(t, b)
	assert.NoError(t, err)

	// Token correct
	b, err = ms.Verify(nil, "token", "uid")
	assert.True(t, b)
	assert.NoError(t, err)
}
