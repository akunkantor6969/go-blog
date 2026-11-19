package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		KirimError(w, "Method tidak diizinkan", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		KirimError(w, "Format JSON tidak valid", http.StatusBadRequest)
		return
	}

	expectedUser := os.Getenv("API_USER")
	expectedPass := os.Getenv("API_PASS")
	if expectedUser == "" {
		expectedUser = "admin"
	}
	if expectedPass == "" {
		expectedPass = "rahasia123"
	}

	if req.Username != expectedUser || req.Password != expectedPass {
		KirimError(w, "Username atau password salah", http.StatusUnauthorized)
		return
	}

	secret := []byte(os.Getenv("JWT_SECRET"))
	claims := jwt.MapClaims{
		"sub": req.Username,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secret)
	if err != nil {
		KirimError(w, "Gagal membuat token", http.StatusInternalServerError)
		return
	}

	KirimSukses(w, "Login berhasil", map[string]any{
		"token":      signedToken,
		"token_type": "Bearer",
		"expires_in": 86400,
	}, http.StatusOK)
}
