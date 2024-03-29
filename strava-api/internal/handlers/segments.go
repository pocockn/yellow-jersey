package handlers

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"

	"yellow-jersey/pkg/logs"
)

// StarredSegments returns a list of starred segments associated with a Strava user.
// TODO: If we receive an error that we aren't authorised for Strava, create new access token.
func (h *Handlers) StarredSegments(c echo.Context) error {
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

	segs, err := h.strava.GetStarredSegments(u.AccessToken)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"segments": segs,
	})
}

// DetailedSegmentsFromStarredSegments fetches a list of detailed segments from Strava based on a users starred segments.
// We need the detailed view, so we get access to the Poly-lines of each segment. A polyline allows us to visualise the
// segment on a map.
// TODO: If we receive an error that we aren't authorised for Strava, create new access token.
func (h *Handlers) DetailedSegmentsFromStarredSegments(c echo.Context) error {
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

	segs, err := h.strava.GetStarredSegments(u.AccessToken)
	if err != nil {
		return err
	}

	var segmentIDs []int
	for _, seg := range segs {
		segmentIDs = append(segmentIDs, seg.Id)
	}

	detailedSegments, err := h.strava.GetSegments(u.AccessToken, segmentIDs)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"segments": detailedSegments,
	})
}

// DetailedSegments fetches segments based off a list of IDs.
func (h *Handlers) DetailedSegments(c echo.Context) error {
	// TODO: Add a helper method for fetching JWT tokens from request
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := claims["sub"].(string)

	logs.Logger.Info().Msgf("fetching user %s", id)
	u, err := h.user.FetchUser(id)
	if err != nil {
		return err
	}

	eventID := c.Param("event_id")
	if eventID == "" {
		return fmt.Errorf("event id can't be empty")
	}

	evt, err := h.events.FetchEvent(eventID)
	if err != nil {
		return err
	}

	detailedSegments, err := h.strava.GetSegments(u.AccessToken, evt.SegmentIDs)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"segments": detailedSegments,
	})
}
