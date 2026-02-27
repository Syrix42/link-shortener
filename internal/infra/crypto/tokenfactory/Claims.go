package crypto

import "github.com/golang-jwt/jwt/v5"

type AccessClaims struct {
	Role      string `json:"role"`
	SessionID string `json:"sid"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	SessionID string `json:"sid"`
	jwt.RegisteredClaims
}
