package token

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrTokenExpired = fmt.Errorf("token expired")
	ErrTokenInvalid = fmt.Errorf("token invalid")
)

type Payload struct {
	ID        uuid.UUID `json:"id"` // we're using an ID to invalidate tokens in case they're leaked
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// NewPayload creates a new payload for the given username and duration
func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	now := time.Now()

	return &Payload{
		ID:        tokenID,
		Username:  username,
		IssuedAt:  now,
		ExpiredAt: now.Add(duration),
	}, nil
}

func (p *Payload) Valid() error {
	if p.ExpiredAt.Before(time.Now()) {
		return ErrTokenExpired
	}

	return nil
}
