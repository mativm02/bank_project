package token

import (
	"time"

	"github.com/xtgo/uuid"
)

type Payload struct {
	ID        uuid.UUID `json:"id"` // we're using an ID to invalidate tokens in case they're leaked
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}
