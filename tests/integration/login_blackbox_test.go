package integration_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"
)

// ---- Login-specific helpers (prefixed to avoid clashes) ----

func loginBaseURL(t *testing.T) string {
	t.Helper()
	u := os.Getenv("BASE_URL")
	if u == "" {
		t.Fatalf("BASE_URL is not set")
	}
	return u
}

func loginPath() string {
	if p := os.Getenv("LOGIN_PATH"); p != "" {
		return p
	}
	return "/api/v1/auth/login"
}

func loginPostJSON(t *testing.T, url string, body []byte) (*http.Response, []byte) {
	t.Helper()

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	b, _ := io.ReadAll(res.Body)
	_ = res.Body.Close()
	return res, b
}

func loginToJSON(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}

// ---- Response DTO ----

type LoginResponse struct {
	Status          string `json:"status"`
	JWTRefreshToken string `json:"jwtRefreshToken"`
	JWTAccessToken  string `json:"jwtAccessToken"`
}

// ---- Tests ----

func TestLogin_InvalidEmail_Returns400(t *testing.T) {
	url := loginBaseURL(t) + loginPath()
	body := []byte(`{"email":"not-an-email","password":"pw123"}`)

	res, b := loginPostJSON(t, url, body)

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", res.StatusCode, string(b))
	}
}

func TestLogin_Success_Returns200_AndTokens(t *testing.T) {
	email := os.Getenv("LOGIN_EMAIL")
	pass := os.Getenv("LOGIN_PASSWORD")
	if email == "" || pass == "" {
		t.Fatalf("LOGIN_EMAIL and LOGIN_PASSWORD must be set to an existing user")
	}

	url := loginBaseURL(t) + loginPath()
	body := []byte(`{"email":` + loginToJSON(email) + `,"password":` + loginToJSON(pass) + `}`)

	res, b := loginPostJSON(t, url, body)

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", res.StatusCode, string(b))
	}

	var lr LoginResponse
	if err := json.Unmarshal(b, &lr); err != nil {
		t.Fatalf("unmarshal LoginResponse: %v, body=%s", err, string(b))
	}

	// adjust this if your API uses something like "success"
	if lr.Status == "" {
		t.Fatalf("expected non-empty status, got body=%s", string(b))
	}
	if lr.JWTAccessToken == "" {
		t.Fatalf("expected non-empty jwtAccessToken, got body=%s", string(b))
	}
	if lr.JWTRefreshToken == "" {
		t.Fatalf("expected non-empty jwtRefreshToken, got body=%s", string(b))
	}
}

func TestLogin_WrongPassword_ReturnsUnauthorized(t *testing.T) {
	email := os.Getenv("LOGIN_EMAIL")
	pass := os.Getenv("LOGIN_PASSWORD")
	if email == "" || pass == "" {
		t.Fatalf("LOGIN_EMAIL and LOGIN_PASSWORD must be set to an existing user")
	}

	url := loginBaseURL(t) + loginPath()
	body := []byte(`{"email":` + loginToJSON(email) + `,"password":"definitely-wrong"}`)

	res, b := loginPostJSON(t, url, body)

	// prefer 401, but allow common variants depending on your handler design
	if res.StatusCode != http.StatusUnauthorized && res.StatusCode != http.StatusBadRequest && res.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 401/400/403, got %d, body=%s", res.StatusCode, string(b))
	}
}

func TestLogin_UserNotFound_Returns404orUnauthorized(t *testing.T) {
	url := loginBaseURL(t) + loginPath()
	body := []byte(`{"email":"missing-user@example.com","password":"pw123"}`)

	res, b := loginPostJSON(t, url, body)

	// some APIs return 401 to prevent account enumeration
	if res.StatusCode != http.StatusNotFound && res.StatusCode != http.StatusUnauthorized && res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 404/401/400, got %d, body=%s", res.StatusCode, string(b))
	}
}
