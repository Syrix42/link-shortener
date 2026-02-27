package repositories

import (
	"context"

	"github.com/Syrix42/link-shortener/internal/domain"
)

type UserRepository interface {
	Save(ctx context.Context, u domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	//GetById(ctx context.Context , id string)(*entities.User , error)
}

// Commented behaviors will be added on later when more Apis were going to add
