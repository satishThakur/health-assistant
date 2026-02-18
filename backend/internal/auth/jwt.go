package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenService handles JWT generation and validation
type TokenService struct {
	secret        []byte
	tokenDuration time.Duration
}

// NewTokenService creates a new TokenService.
// Returns an error if secret is empty or the default placeholder.
func NewTokenService(secret string, hours int) (*TokenService, error) {
	if secret == "" || secret == "change-me-in-production" {
		return nil, errors.New("JWT_SECRET must be set to a non-default value")
	}
	if len(secret) < 32 {
		return nil, fmt.Errorf("JWT_SECRET must be at least 32 characters, got %d", len(secret))
	}
	return &TokenService{
		secret:        []byte(secret),
		tokenDuration: time.Duration(hours) * time.Hour,
	}, nil
}

type jwtClaims struct {
	jwt.RegisteredClaims
}

// GenerateToken creates a signed HS256 JWT for the given userID.
func (s *TokenService) GenerateToken(userID string) (string, error) {
	now := time.Now()
	c := jwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			Issuer:    "health-assistant",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.tokenDuration)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(s.secret)
}

// ValidateToken parses and validates the token string.
// Returns the userID (subject) on success.
func (s *TokenService) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwtClaims{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return s.secret, nil
		},
		jwt.WithIssuer("health-assistant"),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return "", fmt.Errorf("invalid token: %w", err)
	}

	c, ok := token.Claims.(*jwtClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token claims")
	}

	return c.Subject, nil
}
