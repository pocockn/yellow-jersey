package event

import (
	"time"

	"github.com/google/uuid"
)

// Event - An event is a list of user IDs & Strava segments.
// Event is the core object within the application, users can create events, other users to them and add
// Strava segments to them. When events complete the segments the total times get totaled, so we can calculate
// leaderboards.
type Event struct {
	ID         string   `json:"id" bson:"id"`
	Owner      string   `json:"owner" bson:"owner"`
	Name       string   `json:"name" bson:"name"`
	SegmentIDs []string `json:"segment_ids" bson:"segment_ids"`
	Users      []string `json:"users" bson:"users"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// NewEvent creates an empty event and sets the owner. The owner is the user who created the event.
func NewEvent(owner, name string) *Event {
	return &Event{
		ID:        uuid.New().String(),
		Owner:     owner,
		Name:      name,
		CreatedAt: time.Now().UTC(),
	}
}
