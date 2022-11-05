package user_test

import (
	"testing"

	"yellow-jersey/internal/user"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	accessToken := "accessToken"
	refreshToken := "refreshToken"
	stravaID := "stravaID"

	usr := user.NewUser(accessToken, refreshToken, stravaID)

	assert.Equal(t, accessToken, usr.AccessToken)
	assert.Equal(t, refreshToken, usr.RefreshToken)
	assert.Equal(t, stravaID, usr.StravaID)
	assert.NotEmpty(t, usr.ID)
}
