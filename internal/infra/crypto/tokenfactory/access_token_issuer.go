package crypto

import (
	"context"
	"crypto/rsa"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func IssueAccessToken(ctx context.Context, Role, SessionId, UserId string, privateKey *rsa.PrivateKey) (string, error) {
	now := time.Now().UTC()

	claims := AccessClaims{
		Role:      Role,
		SessionID: SessionId,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   UserId,
			Issuer:    "aAuth",
			Audience:  []string{"my-api"},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
		},
	}
	accesstoken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return accesstoken.SignedString(privateKey)
}
