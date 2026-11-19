package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoginHandler_InvalidCredentials(t *testing.T) {
	t.Setenv("API_USER", "admin")
	t.Setenv("API_PASS", "rahasia123")
	t.Setenv("JWT_SECRET", "test-secret")

	body := `{"username":"admin","password":"salah"}`
	req := httptest.NewRequest(http.MethodPost, "/api/login", strings.NewReader(body))
	rec := httptest.NewRecorder()

	LoginHandler(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

func TestLoginHandler_Success(t *testing.T) {
	t.Setenv("API_USER", "admin")
	t.Setenv("API_PASS", "rahasia123")
	t.Setenv("JWT_SECRET", "test-secret")

	body := `{"username":"admin","password":"rahasia123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/login", strings.NewReader(body))
	rec := httptest.NewRecorder()

	LoginHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var resp map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("gagal parse JSON: %v", err)
	}

	if resp["status"] != true {
		t.Errorf("expected status=true, got %v", resp["status"])
	}
}
