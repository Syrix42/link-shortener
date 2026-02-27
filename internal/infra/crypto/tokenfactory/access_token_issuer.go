package crypto

import (
	"context"
	"crypto/rsa"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func IssueAccessToken(ctx context.Context, IsAdmin bool, SessionId, UserId string, privateKey *rsa.PrivateKey) (string, error) {
	now := time.Now().UTC()

	claims := AccessClaims{
		IsAdmin:   IsAdmin,
		SessionID: SessionId,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   UserId,
			Issuer:    "link-Auth",
			Audience:  []string{"link-resources"},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
		},
	}
	accesstoken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return accesstoken.SignedString(privateKey)
}
