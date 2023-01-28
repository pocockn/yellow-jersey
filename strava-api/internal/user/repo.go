package user

import "errors"

//go:generate mockgen -source=$GOFILE -destination=../../mocks/mock_$GOFILE -package=mocks

// ErrUserNotFound is returned when a user is not found.
var ErrUserNotFound = errors.New("the user was not found in the repository")

// Repository holds methods related to User database actions.
type Repository interface {
	CreateUser(accessToken, refreshToken, stravaID string) (*User, error)
	FetchUser(id string) (*User, error)
	FetchUserByStravaID(stravaID string) (*User, error)
	UpdateUser(user *User) error
}
