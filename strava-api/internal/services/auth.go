package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"

	"yellow-jersey/internal/strava"
)

// Authorise  performs the second part of the OAuth exchange. The client has already been redirected to the
// Strava authorization page, has granted authorization to the application and has been redirected back to the
// defined URL. The code param was returned as a query string param in to the redirect_url.
func (s *Strava) Authorise(code string) (*strava.AuthorizationResponse, error) {
	if code == "" {
		return nil, fmt.Errorf("unable to perform second part of OAuth exchange without token")
	}

	resp, err := s.httpClient.PostForm(
		baseURL+"/oauth/token",
		url.Values{
			"client_id":     {fmt.Sprintf("%d", s.clientID)},
			"client_secret": {s.clientSecret},
			"code":          {code},
		},
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		var response strava.Error
		contents, _ := io.ReadAll(resp.Body)
		if err := json.Unmarshal(contents, &response); err != nil {
			return nil, fmt.Errorf("unable to unmarshal response into JSON %w", err)
		}

		return nil, fmt.Errorf("problem performing auth %w", response)
	}

	var response strava.AuthorizationResponse
	contents, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(contents, &response)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal %+v into auth response", string(contents))
	}

	return &response, nil
}

// AuthorizationURL constructs the url a user should use to authorize this specific application.
func (s *Strava) AuthorizationURL(state string) string {
	u, err := url.Parse(fmt.Sprintf("%s/oauth/authorize", baseURL))
	if err != nil {
		return ""
	}

	queryParams := url.Values{}
	queryParams.Add("client_id", fmt.Sprintf("%d", s.clientID))
	queryParams.Add("response_type", "code")
	queryParams.Add("redirect_uri", callbackURL)
	queryParams.Add("scope", "profile:read_all,activity:read_all")
	queryParams.Add("state", state)
	queryParams.Add("approval_prompt", "force")

	u.RawQuery = queryParams.Encode()
	return u.String()
}
