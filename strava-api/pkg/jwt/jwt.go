package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Manager struct keeps details for JWT token generation.
type Manager struct {
	expires int
	secret  string
}

// New creates a new JWT struct.
func New(expires int, secret string) *Manager {
	return &Manager{
		expires: expires,
		secret:  secret,
	}
}

// Generate generates a token.
func (jm Manager) Generate(id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24 * time.Duration(jm.expires)).Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   id,
	})

	return token.SignedString([]byte(jm.secret))
}

// Parse parses jwt token.
func (jm Manager) Parse(t string) (jwt.MapClaims, bool) {
	var (
		err   error
		token *jwt.Token
	)

	if token, err = jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return false, errors.New("unexpected signing method")
		}
		return []byte(jm.secret), nil
	}); err != nil {
		return nil, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	}

	return nil, false
}
