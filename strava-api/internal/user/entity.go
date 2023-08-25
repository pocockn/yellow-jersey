package user

import (
	"time"

	"github.com/google/uuid"

	"yellow-jersey/internal/strava"
)

// User is our model for a user on the system.
// TODO: Make AccessToken and RefreshToken private so they don't get passed through within JSON responses.
type User struct {
	AccessToken  string `json:"access_token" bson:"access_token"`
	ID           string `json:"id" bson:"_id"`
	RefreshToken string `json:"refresh_token" bson:"refresh_token"`
	StravaID     string `json:"strava_id" bson:"strava_id"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`

	AthleteDetailed strava.AthleteDetailed `json:"athlete_detailed" bson:"athlete_detailed"`
}

// NewUser returns us a new populated user struct. The access_token, refresh_token and strava_id are returned when a
// user successfully completes the oauth2 authentication flow with Strava.
func NewUser(accessToken, refreshToken, stravaID string, ath strava.AthleteDetailed) *User {
	return &User{
		AccessToken:     accessToken,
		RefreshToken:    refreshToken,
		StravaID:        stravaID,
		ID:              uuid.New().String(),
		AthleteDetailed: ath,
	}
}
