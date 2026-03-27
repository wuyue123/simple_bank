package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("expired token")
)

type TokenType byte

const (
	TokenTypeAccessToken  TokenType = 1
	TokenTypeRefreshToken TokenType = 2
)

// Payload is the payload of the token.
// It contains the username of the user and the time the token was issued.
type Payload struct {
	Type     TokenType `json:"token_type"`
	Username string    `json:"username"`
	Role     string    `json:"role"`
	jwt.RegisteredClaims
}

// NewPayload creates a new payload.
// It returns a new payload and an error if the payload cannot be created.
func NewPayload(username string, role string, duration time.Duration, tokenType TokenType) (*Payload, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		Type:     tokenType,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenId.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return payload, nil
}

func (payload *Payload) Valid(tokenType TokenType) error {
	if payload.Type != tokenType {
		return ErrInvalidToken
	}
	if payload.ExpiresAt.Before(time.Now()) {
		return ErrExpiredToken
	}
	return nil
}
