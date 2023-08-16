package event_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"yellow-jersey/internal/event"
	"yellow-jersey/internal/strava"
	"yellow-jersey/internal/user"
)

func TestNewEvent(t *testing.T) {
	evt, err := event.NewEvent("nick", "Girona 2023", time.Now(), time.Now().Add(24*time.Hour))
	assert.NoError(t, err)
	assert.Equal(t, "nick", evt.Owner)
	assert.Equal(t, "Girona 2023", evt.Name)
	assert.NotZero(t, evt.CreatedAt)
	assert.Len(t, evt.SegmentEfforts, 0)
}

func TestNewEvent_Validation(t *testing.T) {
	t.Run("Event must have a name supplied", func(t *testing.T) {
		_, err := event.NewEvent("", "Girona 2023", time.Now(), time.Now().Add(24*time.Hour))
		assert.EqualError(t, err, event.ErrInvalidEvent.Error())
	})

	t.Run("Event must have a owner supplied", func(t *testing.T) {
		_, err := event.NewEvent("Nick", "", time.Now(), time.Now().Add(24*time.Hour))
		assert.EqualError(t, err, event.ErrInvalidEvent.Error())
	})
}

func TestEvent_CreateSegmentEfforts(t *testing.T) {
	t.Run("Segments and users are added to event", func(t *testing.T) {
		evt, err := event.NewEvent("Nick", "Girona 2023", time.Now(), time.Now().Add(24*time.Hour))
		assert.NoError(t, err)

		SegmentIDs := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		evt.Users = []user.User{{ID: "nick"}, {ID: "kirsty"}, {ID: "elliot"}, {ID: "knomey"}}

		for _, segmentID := range SegmentIDs {
			require.NoError(t, evt.AddSegment(segmentID))
		}

		assert.Len(t, evt.SegmentIDs, 10)
		assert.Len(t, evt.Users, 4)
	})
}

func TestEvent_PopulateSegmentEfforts(t *testing.T) {
	evt, err := event.NewEvent("Nick", "Girona 2023", time.Now(), time.Now().Add(24*time.Hour))
	assert.NoError(t, err)

	segmentIDs := []int{1}
	evt.Users = []user.User{{ID: "nick"}, {ID: "kirsty"}, {ID: "elliot"}, {ID: "knomey"}}
	for _, segmentID := range segmentIDs {
		if err := evt.AddSegment(segmentID); err != nil {
			t.Fail()
		}
	}

	t.Run("Segment time is populated", func(t *testing.T) {
		err = evt.PopulateSegmentEfforts([]strava.SegmentEffortDetailed{
			{
				SegmentEffortSummary: strava.SegmentEffortSummary{
					Segment:       strava.SegmentSummary{Id: 34},
					EffortSummary: strava.EffortSummary{ElapsedTime: 3000},
				},
			},
		}, "nick")
		assert.NoError(t, err)
		assert.Equal(t, 3000, evt.SegmentEfforts["nick"][0].EffortSummary.ElapsedTime)
	})
}
