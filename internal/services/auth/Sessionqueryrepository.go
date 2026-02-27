package auth

import "context"

type SessionQueryRepository interface {
	CountSessionsByUserID(ctx context.Context, userID string) (int, error)
}
