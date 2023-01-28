package event_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"yellow-jersey/internal/event"
)

func TestNewEvent(t *testing.T) {
	evt := event.NewEvent("nick", "Girona 2023")
	assert.Equal(t, "nick", evt.Owner)
	assert.Equal(t, "Girona 2023", evt.Name)
	assert.NotZero(t, evt.CreatedAt)
}
