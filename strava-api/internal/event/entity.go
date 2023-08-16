package event

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"yellow-jersey/internal/strava"
	"yellow-jersey/internal/user"
	"yellow-jersey/pkg/logs"
)

var (
	// ErrInvalidEvent is returned if the new event is created with an empty owner or name.
	ErrInvalidEvent = fmt.Errorf("an event must have a valid owner, name, start and finish date")
)

// Event - An event is a list of user IDs & Strava segments.
// Event is the core object within the application, users can create events, add other users to them and add
// Strava segments to them. When users complete segments the times get totaled, and we can calculate
// leaderboards.
// TODO: Add struct validation
type Event struct {
	ID         string      `json:"id" bson:"_id"`
	Owner      string      `json:"owner" bson:"owner"`
	Name       string      `json:"name" bson:"name"`
	SegmentIDs []int       `json:"segment_ids" bson:"segment_ids"`
	Users      []user.User `json:"users" bson:"users"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`

	StartDate  time.Time `json:"start_date" bson:"start_date"`
	FinishDate time.Time `json:"finish_date" bson:"finish_date"`

	SegmentEfforts map[string][]strava.SegmentEffortDetailed `json:"segment_efforts" bson:"segment_efforts"`
}

// NewEvent creates an empty event and sets the owner. The owner is the user who created the event.
func NewEvent(owner, name string, startDate, finishDate time.Time) (*Event, error) {
	// Validate that the event name & event owner is not empty
	if name == "" || owner == "" || startDate.IsZero() || finishDate.IsZero() {
		return &Event{}, ErrInvalidEvent
	}

	return &Event{
		ID:             uuid.New().String(),
		Owner:          owner,
		Name:           name,
		CreatedAt:      time.Now().UTC(),
		SegmentEfforts: make(map[string][]strava.SegmentEffortDetailed),
		StartDate:      startDate,
		FinishDate:     finishDate,
	}, nil
}

// AddSegment adds a segment to an event.
func (e *Event) AddSegment(segID int) error {
	for _, segment := range e.SegmentIDs {
		if segID == segment {
			logs.Logger.Info().Msgf("segment %d already added to event", segID)
			return fmt.Errorf("segment %d already added to event", segID)
		}
	}

	e.SegmentIDs = append(e.SegmentIDs, segID)

	return nil
}

// PopulateSegmentEfforts takes a user's best segment time and sets it on the event.
// Warning, this overwrites any old data stored.
func (e *Event) PopulateSegmentEfforts(segments []strava.SegmentEffortDetailed, userID string) error {
	e.SegmentEfforts[userID] = segments

	return nil
}
