package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"yellow-jersey/pkg/logs"
)

// Register routes for our api.
func (h *Handlers) Register(e *echo.Echo) {
	e.Add(http.MethodGet, "/auth", h.AuthorizationURL)
	e.Add(http.MethodPost, "/exchange_token", h.Authorize)

	// TODO: Pass secret from config
	jwt := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte("secret"),
	})

	authed := e.Group("/user")
	authed.Use(jwt)
	authed.GET("/routes", h.Routes)
	authed.GET("/segments", h.StarredSegments)
	authed.POST("/create-event", h.CreateEvent)
	authed.GET("/events", h.FetchUserEvents)
	authed.GET("/event/:id", h.FetchEvent)

	logs.Logger.Info().Msgf("successfully registered API routes")
}
