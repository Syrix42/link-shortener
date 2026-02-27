package crypto

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

type BcryptComparer struct{}

func NewBcryptComparer() *BcryptComparer {
	return &BcryptComparer{}
}

func (b *BcryptComparer) Compare(ctx context.Context, hashed string, plaintext string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plaintext))
	if err != nil {
		return err
	}
	return nil
}
