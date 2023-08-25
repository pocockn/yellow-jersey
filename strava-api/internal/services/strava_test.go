package services_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"yellow-jersey/internal/services"
	"yellow-jersey/internal/strava"
	"yellow-jersey/mocks"
	"yellow-jersey/testutil"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStrava_GetRoutes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Error response from Strava returns error", func(t *testing.T) {
		stravaErr := strava.ErrorDetailed{
			Resource: "routes",
			Field:    "route",
			Code:     "400",
		}

		b, err := json.Marshal(stravaErr)
		require.NoError(t, err)

		bytes.NewReader(b)

		resp := &http.Response{
			Body: testutil.NoopCloser{Reader: bytes.NewBuffer(b)},
		}

		httpClient := mocks.NewMockHTTPClient(ctrl)
		httpClient.EXPECT().Do(gomock.Any()).Return(resp, nil).Times(1)

		srv := services.NewWithStravaHTTPClient(httpClient)

		routes, err := srv.GetRoutes("id", "access")
		assert.Len(t, routes, 0)
		require.Error(t, err)
	})

	t.Run("Successful route response", func(t *testing.T) {
		// TODO: Replace with proper test fixtures
		routes := []strava.Route{
			{
				Id: 134,
			},
		}

		b, err := json.Marshal(routes)
		require.NoError(t, err)

		bytes.NewReader(b)

		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       testutil.NoopCloser{Reader: bytes.NewBuffer(b)},
		}

		httpClient := mocks.NewMockHTTPClient(ctrl)
		httpClient.EXPECT().Do(gomock.Any()).Return(resp, nil).Times(1)

		srv := services.NewWithStravaHTTPClient(httpClient)

		returnedRoutes, err := srv.GetRoutes("id", "access")
		assert.Len(t, returnedRoutes, 1)
		require.NoError(t, err)
	})

	t.Run("Error performing HTTP request", func(t *testing.T) {
		httpClient := mocks.NewMockHTTPClient(ctrl)
		httpClient.EXPECT().Do(gomock.Any()).Return(nil, fmt.Errorf("big system error")).Times(1)

		srv := services.NewWithStravaHTTPClient(httpClient)

		routes, err := srv.GetRoutes("id", "access")
		assert.Len(t, routes, 0)
		require.Error(t, err)
	})
}

func TestStrava_AuthorizationURL(t *testing.T) {
	srv := services.NewStrava(123, "client_secret")
	url := srv.AuthorizationURL("big_state")
	assert.Equal(
		t,
		"https://www.strava.com/api/v3/oauth/authorize?approval_prompt=force&client_id=123&redirect_uri=http%3A%2F%2Flocalhost%3A3000%2Fcallback&response_type=code&scope=profile%3Aread_all%2Cactivity%3Aread_all&state=big_state",
		url,
	)
}

func TestStrava_Authorise(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Successful authorisation with Strava", func(t *testing.T) {
		authResp := strava.AuthorizationResponse{
			AccessToken:  "access",
			Athlete:      strava.AthleteDetailed{},
			RefreshToken: "refresh",
			State:        "state",
		}

		b, err := json.Marshal(authResp)
		require.NoError(t, err)
		bytes.NewReader(b)

		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       testutil.NoopCloser{Reader: bytes.NewBuffer(b)},
		}

		httpClient := mocks.NewMockHTTPClient(ctrl)
		httpClient.EXPECT().PostForm(
			gomock.Any(), gomock.Any(),
		).Return(resp, nil).Times(1)

		srv := services.NewWithStravaHTTPClient(httpClient)

		authRespReturned, err := srv.Authorise("code")
		require.NoError(t, err)
		assert.Equal(t, "access", authRespReturned.AccessToken)
	})

	t.Run("Unable to authorise user with Strava", func(t *testing.T) {
		stravaErr := strava.Error{
			Message: "We can't authorise you, scum.",
		}

		b, err := json.Marshal(stravaErr)
		require.NoError(t, err)
		bytes.NewReader(b)

		resp := &http.Response{
			StatusCode: http.StatusForbidden,
			Body:       testutil.NoopCloser{Reader: bytes.NewBuffer(b)},
		}

		httpClient := mocks.NewMockHTTPClient(ctrl)
		httpClient.EXPECT().PostForm(
			gomock.Any(), gomock.Any(),
		).Return(resp, nil).Times(1)

		srv := services.NewWithStravaHTTPClient(httpClient)

		_, err = srv.Authorise("code")
		require.Error(t, err)
	})

	t.Run("Error performing HTTP request", func(t *testing.T) {
		httpClient := mocks.NewMockHTTPClient(ctrl)
		httpClient.EXPECT().PostForm(
			gomock.Any(), gomock.Any(),
		).Return(nil, fmt.Errorf("big system error")).Times(1)

		srv := services.NewWithStravaHTTPClient(httpClient)

		_, err := srv.Authorise("code")
		require.Error(t, err)
	})
}

func TestStrava_GetStarredSegments(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Successful segment response", func(t *testing.T) {
		segmentFile, err := os.ReadFile("testdata/segments.json")
		assert.NoError(t, err)

		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       testutil.NoopCloser{Reader: bytes.NewBuffer(segmentFile)},
		}

		httpClient := mocks.NewMockHTTPClient(ctrl)
		httpClient.EXPECT().Do(gomock.Any()).Return(resp, nil).Times(1)

		srv := services.NewWithStravaHTTPClient(httpClient)

		segments, err := srv.GetStarredSegments("access")
		assert.Len(t, segments, 1)
		require.NoError(t, err)
	})
}

func TestStrava_GetDetailedSegments(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Successful detailed segment response", func(t *testing.T) {
		segmentFile, err := os.ReadFile("testdata/segment.json")
		assert.NoError(t, err)

		httpClient := mocks.NewMockHTTPClient(ctrl)
		httpClient.EXPECT().Do(gomock.Any()).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       testutil.NoopCloser{Reader: bytes.NewBuffer(segmentFile)},
		}, nil).Times(1)
		httpClient.EXPECT().Do(gomock.Any()).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       testutil.NoopCloser{Reader: bytes.NewBuffer(segmentFile)},
		}, nil).Times(1)

		srv := services.NewWithStravaHTTPClient(httpClient)

		segments, err := srv.GetSegments("access", []int{1, 2})
		assert.Len(t, segments, 2)
		require.NoError(t, err)
	})

	t.Run("Error from detailed segments", func(t *testing.T) {
		stravaErr := strava.Error{
			Message: "We can't authorise you, scum.",
		}

		b, err := json.Marshal(stravaErr)
		require.NoError(t, err)
		bytes.NewReader(b)

		resp := &http.Response{
			StatusCode: http.StatusForbidden,
			Body:       testutil.NoopCloser{Reader: bytes.NewBuffer(b)},
		}

		httpClient := mocks.NewMockHTTPClient(ctrl)
		httpClient.EXPECT().Do(gomock.Any()).Return(resp, nil).Times(1)

		srv := services.NewWithStravaHTTPClient(httpClient)

		_, err = srv.GetSegments("access", []int{1})
		assert.Error(t, err)
	})
}

func TestStrava_GetSegmentEffort(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Successful segment effort response", func(t *testing.T) {
		segmentFile, err := os.ReadFile("testdata/segment_efforts.json")
		assert.NoError(t, err)

		httpClient := mocks.NewMockHTTPClient(ctrl)
		httpClient.EXPECT().Do(gomock.Any()).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       testutil.NoopCloser{Reader: bytes.NewBuffer(segmentFile)},
		}, nil).Times(1)

		srv := services.NewWithStravaHTTPClient(httpClient)

		segments, err := srv.GetSegmentEfforts("access", []int{788127}, time.Now(), time.Now())
		assert.Len(t, segments, 1)
		require.NoError(t, err)
	})

	t.Run("Error from detailed segments", func(t *testing.T) {
		stravaErr := strava.Error{
			Message: "We can't authorise you, scum.",
		}

		b, err := json.Marshal(stravaErr)
		require.NoError(t, err)
		bytes.NewReader(b)

		resp := &http.Response{
			StatusCode: http.StatusForbidden,
			Body:       testutil.NoopCloser{Reader: bytes.NewBuffer(b)},
		}

		httpClient := mocks.NewMockHTTPClient(ctrl)
		httpClient.EXPECT().Do(gomock.Any()).Return(resp, nil).Times(1)

		srv := services.NewWithStravaHTTPClient(httpClient)

		_, err = srv.GetSegmentEfforts("access", []int{1}, time.Now(), time.Now())
		assert.Error(t, err)
	})
}

func TestStrava_GetAthleteDetailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Successful athlete detailed response", func(t *testing.T) {
		segmentFile, err := os.ReadFile("testdata/athlete_detailed.json")
		assert.NoError(t, err)

		httpClient := mocks.NewMockHTTPClient(ctrl)
		httpClient.EXPECT().Do(gomock.Any()).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       testutil.NoopCloser{Reader: bytes.NewBuffer(segmentFile)},
		}, nil).Times(1)

		srv := services.NewWithStravaHTTPClient(httpClient)

		athlete, err := srv.GetAthleteDetailed(12345, "1234")
		assert.Equal(t, "Marianne", athlete.FirstName)
		require.NoError(t, err)
	})

	t.Run("Error from detailed athlete", func(t *testing.T) {
		stravaErr := strava.Error{
			Message: "We can't authorise you, scum.",
		}

		b, err := json.Marshal(stravaErr)
		require.NoError(t, err)
		bytes.NewReader(b)

		resp := &http.Response{
			StatusCode: http.StatusForbidden,
			Body:       testutil.NoopCloser{Reader: bytes.NewBuffer(b)},
		}

		httpClient := mocks.NewMockHTTPClient(ctrl)
		httpClient.EXPECT().Do(gomock.Any()).Return(resp, nil).Times(1)

		srv := services.NewWithStravaHTTPClient(httpClient)

		_, err = srv.GetAthleteDetailed(1234, "234")
		assert.Error(t, err)
	})
}
