package setup

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"yellow-jersey/config"
	"yellow-jersey/internal/event"
	"yellow-jersey/internal/handlers"
	"yellow-jersey/internal/services"
	"yellow-jersey/internal/user"
	"yellow-jersey/pkg/logs"
	iMid "yellow-jersey/pkg/middleware"
)

// Setup runs all the setup for the Strava-API. This includes database connections and setting up routes.
func Setup() (*echo.Echo, error) {
	logs.New(logs.WithWriter(os.Stderr), logs.WithDebug(), logs.WithService("strava-api"))
	cfg, err := config.NewConfig()
	if err != nil {
		return nil, fmt.Errorf("unable to create config %w", err)
	}

	userMongo, err := user.NewMongoRepository(context.Background(), cfg.DB.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("unable to start Mongo %w", err)
	}

	srv := services.NewStrava(cfg.Strava.ClientID, cfg.Strava.ClientSecret)
	usr, err := services.NewUser(services.WithUserRepository(userMongo))
	if err != nil {
		return nil, err
	}

	evtMongo, err := event.NewMongoRepository(context.Background(), cfg.DB.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("unable to start Mongo %w", err)
	}
	evt, err := services.NewEvent(services.WithEventsRepository(evtMongo))
	if err != nil {
		return nil, err
	}

	// TODO: Hold JWT secret in config
	h := handlers.New(srv, usr, evt, "secret")

	e := echo.New()
	e.Use(
		middleware.RequestID(),
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		}),
		iMid.LoggerZerolog(logs.Logger),
	)
	h.Register(e)

	return e, nil
}
