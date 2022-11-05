package services

import (
	"yellow-jersey/internal/event"
	"yellow-jersey/pkg/logs"
)

// EventsConfig is an alias for a function that will take in a pointer to an Event and modify it
type EventsConfig func(u *Events) error

// Events is an implementation of the Event service.
type Events struct {
	repo event.Repo
}

// NewEvent takes a variable amount of UserConfig functions and returns a new User service
func NewEvent(cfgs ...EventsConfig) (*Events, error) {
	os := &Events{}
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

// WithEventsRepository applies a given user repository to the User service
func WithEventsRepository(er event.Repo) EventsConfig {
	return func(e *Events) error {
		e.repo = er
		return nil
	}
}

// CreateEvent creates a user within our database after a successful oauth2 authentication with Strava.
func (e *Events) CreateEvent(owner, name string) (string, error) {
	evt, err := e.repo.Create(owner, name)
	if err != nil {
		return "", err
	}

	logs.Logger.Info().Msgf("successfully created new evt %s for user %s", evt.ID, owner)

	return evt.ID, nil
}

// FetchEvent fetches a user via our internal ID.
func (e *Events) FetchEvent(id string) (*event.Event, error) {
	evt, err := e.repo.Fetch(id)
	if err != nil {
		return nil, err
	}

	return evt, nil
}

// FetchUserEvents fetches a user via our internal ID.
func (e *Events) FetchUserEvents(id string) ([]*event.Event, error) {
	evts, err := e.repo.FetchUserEvents(id)
	if err != nil {
		return nil, err
	}

	return evts, nil
}

// UpdateEvent updates an event within Mongo.
func (e *Events) UpdateEvent(evt *event.Event) error {
	return e.repo.Update(evt)
}
