package domain

import (
	"errors"
	"time"
)

type Session struct {
	ID          string // sid
	UserID      string
	RefreshHash string
	ExpiresAt   time.Time
	RevokedAt   *time.Time
	CreatedAt   time.Time
	LastUsedAt  *time.Time
}

var ErrTooManyActiveSessions = errors.New("active sessions cannot be more than 5")

func NewSession(ID, UserId, RefreshHash string, ExpiresAt, CreatedAt time.Time) *Session {
	return &Session{
		ID:          ID,
		UserID:      UserId,
		RefreshHash: RefreshHash,
		ExpiresAt:   ExpiresAt,
		CreatedAt:   CreatedAt,
	}
}

func EnsureMaxActiveSessions(activeCount, max int) error {
	if activeCount >= max {
		return ErrTooManyActiveSessions
	}
	return nil
}
