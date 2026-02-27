package auth

import (
	"context"
	"time"

	"net/mail"

	"github.com/Syrix42/link-shortener/internal/entities"
	"github.com/Syrix42/link-shortener/internal/services/repositories"
	"github.com/google/uuid"
)

type RegisterService struct {
	UserRepo repositories.UserRepository
	Hasher   PasswordHasher
}

func NewRegisterService(UserRepo repositories.UserRepository,
	hasher PasswordHasher) *RegisterService {
	return &RegisterService{
		UserRepo: UserRepo,
		Hasher:   hasher,
	}
}

func (r *RegisterService) Register(ctx context.Context, Email, password string) error {

	ValidEmail, err := mail.ParseAddress(Email)
	if err != nil {
		return ErrInvalidEmailFormat
	}

	existing, err := r.UserRepo.GetByEmail(ctx, ValidEmail.Address)
	if existing != nil {
		return ErrUserAlreadyExists
	}
	if err != nil {
		return err
	}
	Hashed, err := r.Hasher.Hash(ctx, password)
	if err != nil {
		return err
	}
	user := entities.NewUser(uuid.NewString(), ValidEmail.Address, Hashed, true, false, time.Now().UTC(), time.Now().UTC())

	r.UserRepo.Save(ctx, *user)
	return nil

}
