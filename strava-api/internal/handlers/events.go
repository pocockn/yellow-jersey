package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"

	"yellow-jersey/internal/event"
	"yellow-jersey/pkg/logs"
)

// CreateEvent creates a new internal event.
func (h *Handlers) CreateEvent(c echo.Context) error {
	id := extractUserIDFromRequest(c)
	logs.Logger.Info().Msgf("creating event for user %s", id)

	var evt event.Event
	err := json.NewDecoder(c.Request().Body).Decode(&evt)
	if err != nil {
		return err
	}

	// TODO: Send down the full event struct rather than individual fields.
	evtID, err := h.events.CreateEvent(id, evt.Name, evt.StartDate, evt.FinishDate)
	if err != nil {
		return err
	}
	logs.Logger.Info().Msgf("created event %s", evtID)

	return nil
}

// FetchEvent fetches a single event from the database.
func (h *Handlers) FetchEvent(c echo.Context) error {
	id := c.Param("event_id")
	evt, err := h.events.FetchEvent(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"event": evt,
	})
}

// FetchUserEvents fetches all events for a user.
func (h *Handlers) FetchUserEvents(c echo.Context) error {
	id := extractUserIDFromRequest(c)
	logs.Logger.Info().Msgf("fetching events for user %s", id)

	evts, err := h.events.FetchUserEvents(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"events": evts,
	})
}

// GetUsers gets all users to add to an event.
func (h *Handlers) GetUsers(c echo.Context) error {
	usrs, err := h.user.FetchAll()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"users": usrs,
	})
}

// AddSegmentToEvent adds a segment to an event.
func (h *Handlers) AddSegmentToEvent(c echo.Context) error {
	id := c.Param("event_id")
	if id == "" {
		return fmt.Errorf("event id can't be empty")
	}

	segmentIDInt, err := pathParamInt(c, "segment_id")
	if err != nil {
		return err
	}

	evt, err := h.events.FetchEvent(id)
	if err != nil {
		return err
	}

	if err := evt.AddSegment(segmentIDInt); err != nil {
		return err
	}

	if err := h.events.UpdateEvent(evt); err != nil {
		return err
	}
	logs.Logger.Info().Msgf("added segment %s to event %+v", id, segmentIDInt)

	return nil
}

// AddUserSegmentTimesToEvent takes all the segments added to an event and finds the users best time on those segments.
// We can then use the segment times to work out who has been the fastest.
// TODO: The segment time should be within a specific time period to ensure it happened during the event.
func (h *Handlers) AddUserSegmentTimesToEvent(c echo.Context) error {
	eventID := c.Param("event_id")
	if eventID == "" {
		return fmt.Errorf("event_id or user_id can't be empty")
	}

	evt, err := h.events.FetchEvent(eventID)
	if err != nil {
		return fmt.Errorf("unable to fetch event %s", eventID)
	}

	// TODO: Extract this out into a method, this piece of code is duplicated in lots of places
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := claims["sub"].(string)

	logs.Logger.Info().Msgf("fetching segments for user %s", id)

	u, err := h.user.FetchUser(id)
	if err != nil {
		return err
	}

	if u == nil {
		return fmt.Errorf("nil user, unable to perform request to Strava")
	}

	segmentEfforts, err := h.strava.GetSegmentEfforts(u.AccessToken, evt.SegmentIDs, evt.StartDate, evt.FinishDate)
	if err != nil {
		return err
	}

	if err := evt.PopulateSegmentEfforts(segmentEfforts, u.ID); err != nil {
		return fmt.Errorf("unable to populate segment efforts on event")
	}

	if err := h.events.UpdateEvent(evt); err != nil {
		return err
	}

	logs.Logger.Debug().Msgf(
		"successfully added %d segments to event %s for user %s", len(segmentEfforts), evt.ID, u.ID,
	)
	return nil
}

// AddUserToEvent adds a user to an event.
func (h *Handlers) AddUserToEvent(c echo.Context) error {
	eventID := c.Param("event_id")
	userID := c.Param("user_id")
	if eventID == "" || userID == "" {
		return fmt.Errorf("event_id or user_id can't be empty")
	}

	evt, err := h.events.FetchEvent(eventID)
	if err != nil {
		return err
	}

	for _, user := range evt.Users {
		if userID == user.ID {
			logs.Logger.Info().Msgf("user %s already added to event", userID)
			return echo.NewHTTPError(http.StatusConflict, "user already added to event")
		}
	}

	user, err := h.user.FetchUser(userID)
	if err != nil {
		return err
	}

	evt.Users = append(evt.Users, *user)

	if err := h.events.UpdateEvent(evt); err != nil {
		return err
	}
	logs.Logger.Info().Msgf("added user %s to event %+v", userID, eventID)

	return nil
}

// UpdateEvent updates an event in our database from the payload.
func (h *Handlers) UpdateEvent(c echo.Context) error {
	e := new(event.Event)
	if err := c.Bind(e); err != nil {
		return err
	}

	if err := h.events.UpdateEvent(e); err != nil {
		return err
	}
	logs.Logger.Info().Msgf("updated event %+v", e)

	return nil
}

func pathParamInt(c echo.Context, name string) (int, error) {
	param := c.Param(name)
	result, err := strconv.Atoi(param)
	if err != nil {
		return 0, err
	}

	return result, nil
}
