package handlers

import (
	"yellow-jersey/internal/services"
	"yellow-jersey/pkg/jwt"
)

// Handlers deals with the incoming requests to the API.
type Handlers struct {
	strava *services.Strava
	events *services.Events
	user   *services.User
	jwt    *jwt.Manager
}

// New creates a new Handlers instance to handle requests to the API.
func New(strava *services.Strava, user *services.User, events *services.Events, jwtSecret string) *Handlers {
	return &Handlers{
		strava: strava,
		user:   user,
		events: events,
		// TODO: Set the expiry to how long Strava allows a token to be used
		jwt: jwt.New(1000, jwtSecret),
	}
}
