package user

import (
	"errors"

	"yellow-jersey/internal/strava"
)

//go:generate mockgen -source=$GOFILE -destination=../../mocks/mock_user_$GOFILE -package=mocks

// ErrUserNotFound is returned when a user is not found.
var ErrUserNotFound = errors.New("the user was not found in the repository")

// Repository holds methods related to User database actions.
type Repository interface {
	CreateUser(accessToken, refreshToken, stravaID string, ath strava.AthleteDetailed) (*User, error)
	FetchAll() ([]*User, error)
	FetchUser(id string) (*User, error)
	FetchUserByStravaID(stravaID string) (*User, error)
	UpdateUser(user *User) error
}
