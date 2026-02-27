package crypto

import (
	"context"
	"crypto/rsa"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func IssueRefreshToken(ctx context.Context, SessionID, Userid string, privateKey *rsa.PrivateKey) (string, error) {
	now := time.Now().UTC()
	claims := RefreshClaims{
		SessionID: SessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   Userid,
			Issuer:    "link-Auth",
			Audience:  []string{"revocation"},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute * 10080)),
		},
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return jwtToken.SignedString(privateKey)

}
