package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"yellow-jersey/internal/event"
	"yellow-jersey/internal/handlers"
	"yellow-jersey/internal/services"
	"yellow-jersey/internal/user"
	"yellow-jersey/mocks"
	"yellow-jersey/testutil"
)

func TestHandlers_UpdateEvent(t *testing.T) {
	eventJSON := `{"name":"Croatia 2024","segment_ids": [1,2,3], "users": [{"ID":"1"},{"ID":"2"},{"ID":"3"}]}`
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	evt := new(event.Event)
	evt.Name = "Croatia 2024"
	evt.SegmentIDs = []int{1, 2, 3}
	evt.Users = []user.User{{ID: "1"}, {ID: "2"}, {ID: "3"}}

	eventMock := mocks.NewMockRepo(ctrl)
	eventMock.EXPECT().Update(evt).
		Return(nil).Times(1)

	eventsSrv, err := services.NewEvent(services.WithEventsRepository(eventMock))
	assert.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/event/1234", strings.NewReader(eventJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := handlers.New(nil, nil, eventsSrv, "secret")
	assert.NoError(t, h.UpdateEvent(c))
}

func TestHandlers_Add_Segment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	evt := new(event.Event)
	evt.Name = "Croatia 2024"
	evt.SegmentIDs = []int{1234}
	evt.Users = []user.User{{ID: "1"}, {ID: "2"}, {ID: "3"}}

	eventMock := mocks.NewMockRepo(ctrl)
	eventMock.EXPECT().Fetch("1234").Return(evt, nil)
	eventMock.EXPECT().Update(evt).
		Return(nil).Times(1)

	eventsSrv, err := services.NewEvent(services.WithEventsRepository(eventMock))
	assert.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/", testutil.NoopCloser{})
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("event/:event_id/segment/:segment_id")
	c.SetParamNames("event_id", "segment_id")
	c.SetParamValues("1234", "12345")

	h := handlers.New(nil, nil, eventsSrv, "secret")
	assert.NoError(t, h.AddSegmentToEvent(c))
}

func TestHandlers_Add_User(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	evt := new(event.Event)
	evt.Name = "Croatia 2024"
	evt.SegmentIDs = []int{1234}
	evt.Users = []user.User{{ID: "1"}, {ID: "2"}, {ID: "3"}}

	eventMock := mocks.NewMockRepo(ctrl)
	eventMock.EXPECT().Fetch("1234").Return(evt, nil)
	eventMock.EXPECT().Update(evt).
		Return(nil).Times(1)

	// TODO: Change mock constructor name to include the type
	userMock := mocks.NewMockRepository(ctrl)
	userMock.EXPECT().FetchUser("12345").Return(&user.User{ID: "12345"}, nil)

	userSrv, err := services.NewUser(services.WithUserRepository(userMock))
	assert.NoError(t, err)
	eventsSrv, err := services.NewEvent(services.WithEventsRepository(eventMock))
	assert.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/", testutil.NoopCloser{})
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("event/:event_id/users/:user_id")
	c.SetParamNames("event_id", "user_id")
	c.SetParamValues("1234", "12345")

	h := handlers.New(nil, userSrv, eventsSrv, "secret")
	assert.NoError(t, h.AddUserToEvent(c))
}

func TestHandlers_Add_User_Already_Added(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	evt := new(event.Event)
	evt.Name = "Croatia 2024"
	evt.SegmentIDs = []int{12345}
	evt.Users = []user.User{{ID: "1"}, {ID: "2"}, {ID: "3"}}

	eventMock := mocks.NewMockRepo(ctrl)
	eventMock.EXPECT().Fetch("1234").Return(evt, nil)

	eventsSrv, err := services.NewEvent(services.WithEventsRepository(eventMock))
	assert.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/", testutil.NoopCloser{})
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("event/:event_id/users/:user_id")
	c.SetParamNames("event_id", "user_id")
	c.SetParamValues("1234", "1")

	h := handlers.New(nil, nil, eventsSrv, "secret")
	assert.Error(t, h.AddUserToEvent(c))
}

func TestHandlers_Add_Segment_Already_Added(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	evt := new(event.Event)
	evt.Name = "Croatia 2024"
	evt.SegmentIDs = []int{12345}
	evt.Users = []user.User{{ID: "1"}, {ID: "2"}, {ID: "3"}}

	eventMock := mocks.NewMockRepo(ctrl)
	eventMock.EXPECT().Fetch("1234").Return(evt, nil)

	eventsSrv, err := services.NewEvent(services.WithEventsRepository(eventMock))
	assert.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/", testutil.NoopCloser{})
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("event/:event_id/segment/:segment_id")
	c.SetParamNames("event_id", "segment_id")
	c.SetParamValues("1234", "12345")

	h := handlers.New(nil, nil, eventsSrv, "secret")
	assert.Error(t, h.AddSegmentToEvent(c))
}

func TestHandlers_GetUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usersMock := mocks.NewMockRepository(ctrl)
	usersMock.EXPECT().FetchAll().
		Return([]*user.User{{ID: "1"}, {ID: "2"}}, nil).Times(1)

	userService, err := services.NewUser(services.WithUserRepository(usersMock))
	assert.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/", testutil.NoopCloser{})
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := handlers.New(nil, userService, nil, "secret")
	assert.NoError(t, h.GetUsers(c))
}
