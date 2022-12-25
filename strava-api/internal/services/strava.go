package services

//go:generate mockgen -source=$GOFILE -destination=../../mocks/mock_$GOFILE -package=mocks

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/sync/errgroup"

	"yellow-jersey/internal/strava"
	"yellow-jersey/pkg/logs"
)

const (
	// TODO: Should come from config
	baseURL = "http://mock-strava-api:8082"
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

// GetSegments returns a list of detailed segments, this gives us access to the Polyline, so we can draw
// the segment on a map for the user.
func (s *Strava) GetSegments(accessToken string, ids []int) ([]strava.Segment, error) {
	var segments []strava.Segment

	var g errgroup.Group
	for _, id := range ids {
		segmentID := id
		g.Go(func() error {
			resp, err := s.run("GET", fmt.Sprintf("/segments/%d", segmentID), nil, accessToken)
			if err != nil {
				return err
			}

			var segment strava.Segment
			if err = json.Unmarshal(resp, &segment); err != nil {
				return err
			}

			segments = append(segments, segment)

			return nil
		})
	}

	return segments, g.Wait()
}

func (s *Strava) run(method, path string, params map[string]interface{}, token string) ([]byte, error) {
	var err error

	values := make(url.Values)
	for k, v := range params {
		values.Set(k, fmt.Sprintf("%v", v))
	}

	logs.Logger.Info().Msgf("performing request %s", baseURL+path)
	var req *http.Request
	switch method {
	case http.MethodPost:
		req, err = http.NewRequest(http.MethodPost, baseURL+path, strings.NewReader(values.Encode()))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	default:
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
