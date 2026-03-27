package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const minSecretKeyLength = 32

type JwtMaker struct {
	secret string
}

func NewJwtMaker(secret string) (Maker, error) {
	if len(secret) < minSecretKeyLength {
		return nil, fmt.Errorf("secret key length must be at least %d", minSecretKeyLength)
	}
	return &JwtMaker{secret: secret}, nil
}

// CreateToken creates a new JWT token.
// It returns a token and an error if the token cannot be created.
// The token type is specified by the tokenType parameter.
// The token will be valid for the duration duration.
func (m *JwtMaker) CreateToken(username string, role string, duration time.Duration, tokenType TokenType) (string, *Payload, error) {
	payload, err := NewPayload(username, role, duration, tokenType)

	if err != nil {
		return "", nil, err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	ss, err := token.SignedString([]byte(m.secret))
	return ss, payload, err
}

// VerifyToken verifies a JWT token.
// It returns a payload and an error if the token cannot be verified.
// The token type is specified by the tokenType parameter.
func (m *JwtMaker) VerifyToken(tokenString string, tokenType TokenType) (*Payload, error) {

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(m.secret), nil
	}

	jwtToken, err := jwt.ParseWithClaims(tokenString, &Payload{}, keyFunc)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	err = payload.Valid(tokenType)
	if err != nil {
		return nil, err
	}

	return payload, nil
}
