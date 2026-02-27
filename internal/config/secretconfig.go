package config

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var loadOnce sync.Once

func loadDotEnvOnce() {
	loadOnce.Do(func() {
		err := godotenv.Load()
		if err != nil {
			err = godotenv.Load("../.env")
			if err != nil {
				err = godotenv.Load("../../.env")
				_ = err
			}
		}
	})
}

func LoadPublicAccessJWTKey() (*rsa.PublicKey, error) {
	loadDotEnvOnce()
	return parseRSAPublicKeyFromPathEnv("PUBLIC_ACCESS_JWT_SECRET_PATH")
}

func LoadPrivateAccessJWTKey() (*rsa.PrivateKey, error) {
	loadDotEnvOnce()
	return parseRSAPrivateKeyFromPathEnv("PRIVATE_ACCESS_JWT_SECRET_PATH")
}

func LoadPublicRefreshJWTKey() (*rsa.PublicKey, error) {
	loadDotEnvOnce()
	return parseRSAPublicKeyFromPathEnv("PUBLIC_REFRESH_JWT_SECRET_PATH")
}

func LoadPrivateRefreshJWTKey() (*rsa.PrivateKey, error) {
	loadDotEnvOnce()
	return parseRSAPrivateKeyFromPathEnv("PRIVATE_REFRESH_JWT_SECRET_PATH")
}

func parseRSAPrivateKeyFromPathEnv(envKey string) (*rsa.PrivateKey, error) {
	path, err := mustGetEnv(envKey)
	if err != nil {
		return nil, err
	}

	b, err := readFileFlexible(path)
	if err != nil {
		return nil, fmt.Errorf("%s: read key file: %w", envKey, err)
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(b)
	if err != nil {
		return nil, fmt.Errorf("%s: parse RSA private key PEM: %w", envKey, err)
	}

	return key, nil
}

func parseRSAPublicKeyFromPathEnv(envKey string) (*rsa.PublicKey, error) {
	path, err := mustGetEnv(envKey)
	if err != nil {
		return nil, err
	}

	b, err := readFileFlexible(path)
	if err != nil {
		return nil, fmt.Errorf("%s: read key file: %w", envKey, err)
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(b)
	if err != nil {
		return nil, fmt.Errorf("%s: parse RSA public key PEM: %w", envKey, err)
	}

	return key, nil
}

func mustGetEnv(key string) (string, error) {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return "", fmt.Errorf("missing env var: %s", key)
	}
	return v, nil
}

func readFileFlexible(path string) ([]byte, error) {
	path = strings.TrimSpace(path)

	if filepath.IsAbs(path) {
		return os.ReadFile(path)
	}

	if b, err := os.ReadFile(path); err == nil {
		return b, nil
	}

	if b, err := os.ReadFile(filepath.Join("..", path)); err == nil {
		return b, nil
	}

	if b, err := os.ReadFile(filepath.Join("..", "..", path)); err == nil {
		return b, nil
	}

	return nil, errors.New("file not found in working dir or parent dirs")
}
