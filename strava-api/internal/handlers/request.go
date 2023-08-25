package handlers

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func extractUserIDFromRequest(c echo.Context) string {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := claims["sub"].(string)
	return id
}
