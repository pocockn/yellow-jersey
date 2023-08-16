package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"yellow-jersey/pkg/logs"
)

// AuthRequest is the request we send to Strava during the OAuth2 flow to retrieve an access token.
type AuthRequest struct {
	Code  string `json:"code"`
	State string `json:"state"`
	Scope string `json:"scope"`
}

// Authorize performs the second part of the OAuth exchange. The client has already been redirected to the
// Strava authorization page, has granted authorization to the application and has been redirected back to the
// defined URL. The code param was returned as a query string param in to the redirect_url.
func (h *Handlers) Authorize(c echo.Context) error {
	if c.Request().FormValue("error") == "access_denied" {
		logs.Logger.Error().Msg("access denied")
		return fmt.Errorf("access denied from Strava")
	}

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return fmt.Errorf("unable to read body %w", err)
	}

	authReq := AuthRequest{}
	if err = json.Unmarshal(body, &authReq); err != nil {
		return fmt.Errorf("unable to unmarshal body %s : %w", string(body), err)
	}

	resp, err := h.strava.Authorise(authReq.Code)
	if err != nil {
		return err
	}

	u, err := h.user.FetchUserByStravaID(strconv.FormatInt(resp.Athlete.Id, 10))
	if err != nil {
		return err
	}
	if u != nil {
		logs.Logger.Info().Msgf("found user %s in system", u.ID)
		u.AccessToken = resp.AccessToken
		u.RefreshToken = resp.RefreshToken
		if err = h.user.UpdateUser(u); err != nil {
			return err
		}

		return h.generateJWT(c, u.ID)
	}

	athlete, err := h.strava.GetAthleteDetailed(int(resp.Athlete.Id), resp.AccessToken)
	if err != nil {
		return err
	}

	userID, err := h.user.CreateUser(
		resp.AccessToken,
		resp.RefreshToken,
		strconv.FormatInt(resp.Athlete.Id, 10),
		athlete,
	)
	if err != nil {
		return err
	}

	return h.generateJWT(c, userID)
}

// AuthorizationURL constructs the url a user should use to authorize this specific application.
func (h *Handlers) AuthorizationURL(c echo.Context) error {
	// TODO: Generate state
	if _, err := c.Response().Writer.Write([]byte(h.strava.AuthorizationURL("state1"))); err != nil {
		return err
	}

	return nil
}

// generateJWT generates a JWT for a user.
func (h *Handlers) generateJWT(c echo.Context, userID string) error {
	token, err := h.jwt.Generate(userID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
}
