package repositories

import (
	"context"

	"github.com/Syrix42/link-shortener/internal/entities"
)

type UserRepository interface {
	Save(ctx context.Context, u entities.User) error
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	//GetById(ctx context.Context , id string)(*entities.User , error)
	//IncrementActiveSessions(userID string) error
	//DecrementActiveSessions(userID string) error
}

// Commented behaviors will be added on later when more Apis were going to add
