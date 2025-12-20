package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// локальный тест-дабл репозитория
type stubRepo struct {
	users map[string]User
}

func (r stubRepo) ByEmail(email string) (User, error) {
	u, ok := r.users[email]
	if !ok {
		return User{}, ErrNotFound
	}
	return u, nil
}

func TestService_FindIDByEmail(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		repo := stubRepo{
			users: map[string]User{
				"user@example.com": {ID: 42, Email: "user@example.com"},
			},
		}
		svc := New(repo)

		id, err := svc.FindIDByEmail("user@example.com")

		require.NoError(t, err)
		assert.Equal(t, int64(42), id)
	})

	t.Run("not found", func(t *testing.T) {
		repo := stubRepo{
			users: map[string]User{},
		}
		svc := New(repo)

		id, err := svc.FindIDByEmail("missing@example.com")

		assert.ErrorIs(t, err, ErrNotFound)
		assert.Equal(t, int64(0), id)
	})
}

func TestService_FindIDByEmail_Extra(t *testing.T) {
	repo := stubRepo{
		users: map[string]User{
			"zero@example.com": {ID: 0, Email: "zero@example.com"},
			"multi@example.com": {ID: 99, Email: "multi@example.com"},
		},
	}
	svc := New(repo)

	t.Run("zero ID", func(t *testing.T) {
		id, err := svc.FindIDByEmail("zero@example.com")
		require.NoError(t, err)
		assert.Equal(t, int64(0), id)
	})

	t.Run("another user", func(t *testing.T) {
		id, err := svc.FindIDByEmail("multi@example.com")
		require.NoError(t, err)
		assert.Equal(t, int64(99), id)
	})

	t.Run("empty email", func(t *testing.T) {
		id, err := svc.FindIDByEmail("")
		assert.ErrorIs(t, err, ErrNotFound)
		assert.Equal(t, int64(0), id)
	})
}
