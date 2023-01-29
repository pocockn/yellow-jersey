package user

import (
	"fmt"

	"gorm.io/gorm"
)

// Ensure MySQLRepository implements the Repository interface
var _ Repository = (*MySQLRepository)(nil)

// MySQLRepository is a MySQL implementation of the user.Repository interface.
type MySQLRepository struct {
	db *gorm.DB
}

// NewMySQLRepository returns a MySQLRepository struct with a database connection.
func NewMySQLRepository(db *gorm.DB) *MySQLRepository {
	return &MySQLRepository{db: db}
}

// CreateUser creates a new user within our database.This happens after a successful oauth2 authentication with
// Strava.
func (m *MySQLRepository) CreateUser(accessToken, refreshToken, stravaID string) (*User, error) {
	user := NewUser(accessToken, refreshToken, stravaID)
	if err := m.db.Create(user).Error; err != nil {
		return nil, fmt.Errorf("problem creating user with Strava ID %s within MySQL DB : %w", stravaID, err)
	}

	return user, nil
}

// FetchUser fetches a user by our internal ID.
func (m *MySQLRepository) FetchUser(id string) (*User, error) {
	var fetchedUser *User
	if err := m.db.Raw("SELECT * FROM users WHERE id = ?", id).Scan(&fetchedUser).Error; err != nil {
		return nil, fmt.Errorf("problem fetching user %s from MySQL DB : %w", id, err)
	}

	return fetchedUser, nil
}

// FetchUserByStravaID fetches a user by a Strava ID.
func (m *MySQLRepository) FetchUserByStravaID(id string) (*User, error) {
	var fetchedUser *User
	if err := m.db.Raw("SELECT * FROM users WHERE strava_id = ?", id).Scan(&fetchedUser).Error; err != nil {
		return nil, fmt.Errorf("problem fetching user with Strava ID %s from MySQL DB : %w", id, err)
	}

	return fetchedUser, nil
}

// FetchAll fetches all users. Needs to be implemented for MySQL.
func (m *MySQLRepository) FetchAll() ([]*User, error) {
	return nil, nil
}

// UpdateUser updates a user within the database.
func (m *MySQLRepository) UpdateUser(u *User) error {
	if err := m.db.Save(u).Error; err != nil {
		return fmt.Errorf("unable to update user %s within MySQL DB : %w", u.ID, err)
	}

	return nil
}
