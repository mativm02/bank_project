package token

import "time"

// Maker is an interface for generating and validating tokens
type Maker interface {
	// CreateToken creates a new token for the given username and duration
	CreateToken(username string, duration time.Duration) (string, error)
	// VerifyToken verifies the given token and returns the payload
	VerifyToken(token string) (*Payload, error)
}
