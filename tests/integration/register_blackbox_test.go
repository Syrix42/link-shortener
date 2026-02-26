package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

func baseURL(t *testing.T) string {
	t.Helper()
	u := os.Getenv("BASE_URL")
	if u == "" {
		t.Fatalf("BASE_URL is not set ")
	}
	return u
}

func registerPath() string {
	if p := os.Getenv("REGISTER_PATH"); p != "" {
		return p
	}
	return "/api/v1/auth/register"
}

func postJSON(t *testing.T, url string, body []byte) (*http.Response, []byte) {
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

func TestRegister_InvalidEmail_Returns400(t *testing.T) {
	url := baseURL(t) + registerPath()
	t.Logf("POST %s", url)

	body := []byte(`{"email":"not-an-email","password":"pw123"}`)
	res, b := postJSON(t, url, body)

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", res.StatusCode, string(b))
	}
}

func TestRegister_Success_Returns200or201(t *testing.T) {
	url := baseURL(t) + registerPath()

	email := fmt.Sprintf("u_%d@example.com", time.Now().UnixNano())
	body := []byte(fmt.Sprintf(`{"email":%q,"password":"pw123"}`, email))

	res, b := postJSON(t, url, body)

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		t.Fatalf("expected 200/201, got %d, body=%s", res.StatusCode, string(b))
	}

	_ = json.Valid(b)
}

func TestRegister_DuplicateEmail_ReturnsConflictOrBadRequest(t *testing.T) {
	url := baseURL(t) + registerPath()

	email := fmt.Sprintf("dup_%d@example.com", time.Now().UnixNano())
	body := []byte(fmt.Sprintf(`{"email":%q,"password":"pw123"}`, email))

	res1, b1 := postJSON(t, url, body)
	if res1.StatusCode != http.StatusOK && res1.StatusCode != http.StatusCreated {
		t.Fatalf("expected 200/201 on first register, got %d, body=%s", res1.StatusCode, string(b1))
	}

	res2, b2 := postJSON(t, url, body)

	if res2.StatusCode != http.StatusConflict && res2.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 409/400 on duplicate, got %d, body=%s", res2.StatusCode, string(b2))
	}
}
