package jwt_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"yellow-jersey/pkg/jwt"
)

func TestManager(t *testing.T) {
	t.Parallel()
	manager := jwt.New(100, "secret")

	token, err := manager.Generate("randomID")
	require.NoError(t, err)

	claims, ok := manager.Parse(token)
	assert.True(t, ok)
	require.NoError(t, claims.Valid())

	_, ok = manager.Parse("incorrectToken")
	assert.False(t, ok)
}
