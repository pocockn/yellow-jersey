package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"yellow-jersey/internal/handlers"
	"yellow-jersey/internal/services"
	"yellow-jersey/internal/strava"
	"yellow-jersey/mocks"
	"yellow-jersey/testutil"
)

func TestHandlers_AuthorizationURL(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", testutil.NoopCloser{})
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/exchange_token")

	h := handlers.New(services.NewStrava(1234, "clientSecret"), nil, nil, "secret")
	assert.NoError(t, h.AuthorizationURL(c))
	assert.Equal(
		t,
		"https://www.strava.com/api/v3/oauth/authorize?approval_prompt=force&client_id=1234&redirect_uri=http%3A%2F%2Flocalhost%3A3000%2Fcallback&response_type=code&scope=profile%3Aread_all%2Cactivity%3Aread_all&state=state1",
		rec.Body.String(),
	)
}

// TODO: Complete test.
func TestHandlers_Authorize_NewUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// TODO: Change mock constructor name to include the type
	userMock := mocks.NewMockRepository(ctrl)
	userMock.EXPECT().FetchUserByStravaID("12345").Return(nil, nil)
	userMock.EXPECT().CreateUser(
		"access", "refresh", "stravaID", strava.AthleteDetailed{FTP: 205}).
		Return("userID", nil)

	userSrv, err := services.NewUser(services.WithUserRepository(userMock))
	assert.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/", testutil.NoopCloser{})
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("event/:event_id/users/:user_id")
	c.SetParamNames("event_id", "user_id")
	c.SetParamValues("1234", "12345")

	h := handlers.New(nil, userSrv, nil, "secret")
	assert.NoError(t, h.AddUserToEvent(c))
}
