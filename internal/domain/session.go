package domain

import "time"

type Session struct {
	ID          string // sid
	UserID      string
	RefreshHash string
	ExpiresAt   time.Time
	RevokedAt   *time.Time
	CreatedAt   time.Time
	LastUsedAt  *time.Time
}

func NewSession(ID, UserId, RefreshHash string, ExpiresAt, CreatedAt time.Time) *Session {
	return &Session{
		ID:          ID,
		UserID:      UserId,
		RefreshHash: RefreshHash,
		ExpiresAt:   ExpiresAt,
		CreatedAt:   CreatedAt,
	}
}
