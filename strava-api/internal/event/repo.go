package event

//go:generate mockgen -source=$GOFILE -destination=../mocks/mock_event_$GOFILE -package=mocks

// Repo holds methods related to Event database actions.
type Repo interface {
	Create(ownerID, name string) (*Event, error)
	Fetch(id string) (*Event, error)
	FetchUserEvents(userID string) ([]*Event, error)
	Update(e *Event) error
}
