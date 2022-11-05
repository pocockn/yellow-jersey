package handlers

import (
	"net/http"

	"yellow-jersey/pkg/logs"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// Routes returns a list of routes associated with a Strava user.
// TODO: If we receive an error that we aren't authorised for Strava, create new access token.
func (h *Handlers) Routes(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := claims["sub"].(string)

	logs.Logger.Info().Msgf("fetching routes for user %s", id)

	u, err := h.user.FetchUser(id)
	if err != nil {
		return err
	}

	routes, err := h.strava.GetRoutes(u.StravaID, u.AccessToken)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"routes": routes,
	})
}
