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
