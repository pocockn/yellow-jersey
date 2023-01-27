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
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := claims["sub"].(string)

	logs.Logger.Info().Msgf("creating event for user %s", id)

	var evt event.Event
	err := json.NewDecoder(c.Request().Body).Decode(&evt)
	if err != nil {
		return err
	}

	evtID, err := h.events.CreateEvent(id, evt.Name)
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
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := claims["sub"].(string)

	logs.Logger.Info().Msgf("fetching events for user %s", id)

	evts, err := h.events.FetchUserEvents(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"events": evts,
	})
}

// AddSegment adds a segment to an event.
func (h *Handlers) AddSegment(c echo.Context) error {
	id := c.Param("event_id")
	if id == "" {
		return fmt.Errorf("event id can't be empty")
	}

	segmentIDInt, err := queryParamInt(c, "segment_id")
	if err != nil {
		return err
	}

	evt, err := h.events.FetchEvent(id)
	if err != nil {
		return err
	}

	for _, segment := range evt.SegmentIDs {
		if segmentIDInt == segment {
			logs.Logger.Info().Msgf("segment %d already added to event", segmentIDInt)
			return echo.NewHTTPError(http.StatusConflict, "segment already added to event")

		}
	}

	evt.SegmentIDs = append(evt.SegmentIDs, segmentIDInt)

	if err := h.events.UpdateEvent(evt); err != nil {
		return err
	}
	logs.Logger.Info().Msgf("added segment %s to event %+v", id, segmentIDInt)

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

func queryParamInt(c echo.Context, name string) (int, error) {
	param := c.QueryParam(name)
	result, err := strconv.Atoi(param)
	if err != nil {
		return 0, err
	}

	return result, nil
}
