package services

//go:generate mockgen -source=$GOFILE -destination=../../mocks/mock_$GOFILE -package=mocks

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"yellow-jersey/internal/strava"
	"yellow-jersey/pkg/logs"
)

const (
	baseURL = "https://www.strava.com/api/v3"
)

var (
	// TODO: Should come from config
	callbackURL = fmt.Sprintf("http://localhost:%d/callback", 3000)
)

// Strava communicates with the Strava API.
type Strava struct {
	clientID     int64
	clientSecret string
	httpClient   HTTPClient
}

// NewStrava returns a new client used to communicate with the Strava API.
func NewStrava(clientID int64, clientSecret string) *Strava {
	return &Strava{
		clientID:     clientID,
		clientSecret: clientSecret,
		httpClient:   &http.Client{},
	}
}

// NewWithStravaHTTPClient creates a new Strava struct with an injected client, used for testing.
func NewWithStravaHTTPClient(httpClient HTTPClient) *Strava {
	return &Strava{
		httpClient: httpClient,
	}
}

// HTTPClient is an interface to define the methods required from any kind of
// HTTP Client that will be used by the Client.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
	PostForm(url string, data url.Values) (resp *http.Response, err error)
}

// GetRoutes returns a list of users routes.
func (s *Strava) GetRoutes(stravaID string, accessToken string) ([]strava.Route, error) {
	resp, err := s.run("GET", fmt.Sprintf("/athletes/%s/routes", stravaID), nil, accessToken)
	if err != nil {
		return nil, err
	}

	var routes []strava.Route
	if err = json.Unmarshal(resp, &routes); err != nil {
		return nil, err
	}

	return routes, nil
}

// GetStarredSegments returns a list of segments a user has starred. The user can then add these segments to
// an event.
func (s *Strava) GetStarredSegments(accessToken string) ([]strava.Segment, error) {
	resp, err := s.run("GET", "/segments/starred", nil, accessToken)
	if err != nil {
		return nil, err
	}

	var segments []strava.Segment
	if err = json.Unmarshal(resp, &segments); err != nil {
		return nil, err
	}

	return segments, nil
}

func (s *Strava) run(method, path string, params map[string]interface{}, token string) ([]byte, error) {
	var err error

	values := make(url.Values)
	for k, v := range params {
		values.Set(k, fmt.Sprintf("%v", v))
	}

	logs.Logger.Info().Msgf("performing request %s", baseURL+path)
	var req *http.Request
	if method == "POST" {
		req, err = http.NewRequest("POST", baseURL+path, strings.NewReader(values.Encode()))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		req, err = http.NewRequest(method, baseURL+path+"?"+values.Encode(), nil)
		if err != nil {
			return nil, err
		}
	}

	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		var response strava.Error
		contents, _ := io.ReadAll(resp.Body)
		if err := json.Unmarshal(contents, &response); err != nil {
			return nil, fmt.Errorf("unable to unmarshal response into JSON %w", err)
		}

		return nil, fmt.Errorf("problem performing request %w", response)
	}

	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
