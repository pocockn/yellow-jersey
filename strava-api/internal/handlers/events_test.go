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
	"yellow-jersey/mocks"
)

func TestHandlers_UpdateEvent(t *testing.T) {
	eventJSON := `{"name":"Croatia 2024","segment_ids": ["1","2","3"], "users": ["1","2","3"]}`
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	evt := new(event.Event)
	evt.Name = "Croatia 2024"
	evt.SegmentIDs = []string{"1", "2", "3"}
	evt.Users = []string{"1", "2", "3"}

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
