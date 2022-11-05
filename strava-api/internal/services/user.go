package services

import (
	"errors"

	"yellow-jersey/internal/user"
	"yellow-jersey/pkg/logs"
)

// UserConfig is an alias for a function that will take in a pointer to an User and modify it
type UserConfig func(u *User) error

// User is a implementation of the User service.
type User struct {
	repo user.Repository
}

// NewUser takes a variable amount of UserConfig functions and returns a new User service
// Each UserConfig will be called in the order they are passed in
func NewUser(cfgs ...UserConfig) (*User, error) {
	os := &User{}
	// Apply all Configurations passed in
	for _, cfg := range cfgs {
		// Pass the service into the configuration function
		err := cfg(os)
		if err != nil {
			return nil, err
		}
	}
	return os, nil
}

// WithUserRepository applies a given user repository to the User service
func WithUserRepository(ur user.Repository) UserConfig {
	// return a function that matches the User alias,
	// You need to return this so that the parent function can take in all the needed parameters
	return func(u *User) error {
		u.repo = ur
		return nil
	}
}

// CreateUser creates a user within our database after a successful oauth2 authentication with Strava.
func (u *User) CreateUser(accessToken, refreshToken, stravaID string) (string, error) {
	newUser, err := u.repo.CreateUser(accessToken, refreshToken, stravaID)
	if err != nil {
		return "", err
	}

	logs.Logger.Info().Msgf("successfully created newUser %s with strava id %s", newUser.ID, stravaID)

	return newUser.ID, nil
}

// FetchUser fetches a user via our internal ID.
func (u *User) FetchUser(id string) (*user.User, error) {
	usr, err := u.repo.FetchUser(id)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			logs.Logger.Warn().Msgf("user %s not found", id)
			return nil, nil
		}

		return nil, err
	}

	return usr, nil
}

// FetchUserByStravaID fetches a user via their Strava ID.
func (u *User) FetchUserByStravaID(id string) (*user.User, error) {
	usr, err := u.repo.FetchUserByStravaID(id)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return usr, nil
}

func (u *User) UpdateUser(usr *user.User) error {
	return u.repo.UpdateUser(usr)
}
