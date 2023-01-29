package services_test

import (
	"testing"

	"yellow-jersey/internal/services"
	"yellow-jersey/internal/user"
	"yellow-jersey/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUser_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userMock := mocks.NewMockRepository(ctrl)
	userMock.EXPECT().CreateUser("accessToken", "refreshToken", "stravaID").
		Return(&user.User{ID: "id"}, nil).Times(1)

	srv, err := services.NewUser(services.WithUserRepository(userMock))
	require.NoError(t, err)

	id, err := srv.CreateUser("accessToken", "refreshToken", "stravaID")
	assert.Equal(t, "id", id)
	require.NoError(t, err)
}

func TestUser_FetchUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userMock := mocks.NewMockRepository(ctrl)
	userMock.EXPECT().FetchUser("internalID").
		Return(&user.User{ID: "id"}, nil).Times(1)

	srv, err := services.NewUser(services.WithUserRepository(userMock))
	require.NoError(t, err)

	usr, err := srv.FetchUser("internalID")
	assert.Equal(t, "id", usr.ID)
	require.NoError(t, err)
}

func TestUser_FetchAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userMock := mocks.NewMockRepository(ctrl)
	userMock.EXPECT().FetchAll().
		Return([]*user.User{{ID: "1"}, {ID: "2"}}, nil).Times(1)

	srv, err := services.NewUser(services.WithUserRepository(userMock))
	require.NoError(t, err)

	usrs, err := srv.FetchAll()
	require.NoError(t, err)
	assert.Len(t, usrs, 2)
}

func TestUser_FetchUserByStravaID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userMock := mocks.NewMockRepository(ctrl)
	userMock.EXPECT().FetchUserByStravaID("stravaID").
		Return(&user.User{ID: "id"}, nil).Times(1)

	srv, err := services.NewUser(services.WithUserRepository(userMock))
	require.NoError(t, err)

	usr, err := srv.FetchUserByStravaID("stravaID")
	assert.Equal(t, "id", usr.ID)
	require.NoError(t, err)
}

func TestUser_UpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := user.User{ID: "id"}

	userMock := mocks.NewMockRepository(ctrl)
	userMock.EXPECT().UpdateUser(&u).
		Return(nil).Times(1)

	srv, err := services.NewUser(services.WithUserRepository(userMock))
	require.NoError(t, err)

	require.NoError(t, srv.UpdateUser(&u))
}
