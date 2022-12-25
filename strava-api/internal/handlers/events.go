package handlers

import (
	"encoding/json"
	"net/http"

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
	id := c.Param("id")
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
