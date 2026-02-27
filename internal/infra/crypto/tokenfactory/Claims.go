package crypto

import "github.com/golang-jwt/jwt/v5"

type AccessClaims struct {
	IsAdmin   bool   `json:"is_admin"`
	SessionID string `json:"sid"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	SessionID string `json:"sid"`
	jwt.RegisteredClaims
}
