package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

var loadOnce sync.Once

func loadDotEnvOnce() {
	loadOnce.Do(func() {
		_ = godotenv.Load()
	})
}

func LoadPublicAccessJWTKey() (*rsa.PublicKey, error) {
	loadDotEnvOnce()
	return parseRSAPublicKeyFromEnv("Public_ACESS_JWT_SECRET")
}

func LoadPrivateAccessJWTKey() (*rsa.PrivateKey, error) {
	loadDotEnvOnce()
	return parseRSAPrivateKeyFromEnv("PRIVATE_ACCESS_JWT_SECRET")
}

func LoadPublicRefreshJWTKey() (*rsa.PublicKey, error) {
	loadDotEnvOnce()
	return parseRSAPublicKeyFromEnv("PUBLIC_REFRESH_JWT_SECRET")
}

func LoadPrivateRefreshJWTKey() (*rsa.PrivateKey, error) {
	loadDotEnvOnce()
	return parseRSAPrivateKeyFromEnv("PRIVATE_REFRESH_JWT_SECRET")
}

func parseRSAPublicKeyFromEnv(envKey string) (*rsa.PublicKey, error) {
	raw, err := readEnv(envKey)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode([]byte(raw))
	if block == nil {
		return nil, fmt.Errorf("%s: invalid PEM (no block found)", envKey)
	}

	if block.Type == "RSA PUBLIC KEY" {
		pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("%s: parse PKCS1 public key: %w", envKey, err)
		}
		return pub, nil
	}

	pubAny, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {

		if pub, err2 := x509.ParsePKCS1PublicKey(block.Bytes); err2 == nil {
			return pub, nil
		}
		return nil, fmt.Errorf("%s: parse PKIX public key: %w", envKey, err)
	}

	pub, ok := pubAny.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("%s: not an RSA public key", envKey)
	}
	return pub, nil
}

func parseRSAPrivateKeyFromEnv(envKey string) (*rsa.PrivateKey, error) {
	raw, err := readEnv(envKey)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode([]byte(raw))
	if block == nil {
		return nil, fmt.Errorf("%s: invalid PEM (no block found)", envKey)
	}

	if block.Type == "RSA PRIVATE KEY" {
		priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("%s: parse PKCS1 private key: %w", envKey, err)
		}
		return priv, nil
	}

	privAny, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		// Fallback if it is actually PKCS#1.
		if priv, err2 := x509.ParsePKCS1PrivateKey(block.Bytes); err2 == nil {
			return priv, nil
		}
		return nil, fmt.Errorf("%s: parse PKCS8 private key: %w", envKey, err)
	}

	priv, ok := privAny.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("%s: not an RSA private key", envKey)
	}
	return priv, nil
}

func readEnv(key string) (string, error) {
	val := strings.TrimSpace(os.Getenv(key))
	if val == "" {
		return "", fmt.Errorf("missing env var: %s", key)
	}

	val = strings.ReplaceAll(val, `\n`, "\n")

	if !strings.Contains(val, "BEGIN") {
		return "", errors.New("key does not look like PEM; expected a PEM-encoded RSA key")
	}
	return val, nil
}
