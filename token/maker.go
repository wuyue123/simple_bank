package token

import "time"

// Maker is the interface that defines the methods for creating and verifying tokens.
// It returns a token and an error if the token cannot be created or verified.
type Maker interface {
	CreateToken(username string, role string, duration time.Duration, tokenType TokenType) (string, *Payload, error)
	VerifyToken(token string, tokenType TokenType) (*Payload, error)
}
