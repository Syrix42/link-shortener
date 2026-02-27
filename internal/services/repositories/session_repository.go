package repositories

import (
	"context"

	"github.com/Syrix42/link-shortener/internal/domain"
)

type SessionRepository interface {
	CreateSession(ctx context.Context, Session *domain.Session) error
	RotateRefreshToken(ctx context.Context, Session *domain.Session) error
	DeleteSession(ctx context.Context, SessionID string) error
}

// For following CQRS Pattern the Reads from Session Table will not be here
