package auth

import "context"

type PasswordHasher interface {
	Hash(ctx context.Context, Password string) (string, error)
}
