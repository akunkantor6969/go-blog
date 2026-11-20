package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"go-blog/config"
	"go-blog/handler"

	"github.com/golang-jwt/jwt/v5"

	_ "go-blog/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

// Middleware untuk mencatat (logging) setiap request yang masuk
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Lanjutkan request ke handler berikutnya
		next.ServeHTTP(w, r)

		// Cetak informasi setelah request selesai diproses
		log.Printf("[%s] %s | Durasi: %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// func BasicAuthMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// Izinkan GET tanpa auth (opsional, hapus blok ini jika semua endpoint ingin dikunci)
// 		if r.Method == http.MethodGet {
// 			next.ServeHTTP(w, r)
// 			return
// 		}

// 		expectedUser := os.Getenv("API_USER")
// 		expectedPass := os.Getenv("API_PASS")

// 		user, pass, ok := r.BasicAuth()
// 		userMatch := subtle.ConstantTimeCompare([]byte(user), []byte(expectedUser)) == 1
// 		passMatch := subtle.ConstantTimeCompare([]byte(pass), []byte(expectedPass)) == 1

// 		if !ok || !userMatch || !passMatch {
// 			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
// 			handler.KirimError(w, "Akses ditolak: kredensial tidak valid", http.StatusUnauthorized)
// 			return
// 		}

// 		next.ServeHTTP(w, r)
// 	})
// }

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Biarkan GET dan /api/login bisa diakses tanpa token
		if r.Method == http.MethodGet || r.URL.Path == "/api/login" {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			handler.KirimError(w, "Token Bearer diperlukan", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		secret := []byte(os.Getenv("JWT_SECRET"))

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("metode signing tidak valid")
			}
			return secret, nil
		})

		if err != nil || !token.Valid {
			handler.KirimError(w, "Token tidak valid atau sudah kedaluwarsa", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// @title           go-blog REST API
// @version         1.0
// @description     Dokumentasi REST API untuk manajemen artikel dengan autentikasi JWT.
// @host            localhost:8080
// @BasePath        /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Masukkan token JWT dengan format: Bearer <token>

func main() {
	// Panggil inisialisasi database saat aplikasi pertama kali menyala
	config.InitDatabase()

	if config.DB != nil {
		defer func() {
			if err := config.DB.Close(); err != nil {
				log.Printf("Gagal menutup koneksi database: %v", err)
			} else {
				log.Println("Koneksi database ditutup dengan bersih.")
			}
		}()
	}

	// Daftarkan handler seperti biasa
	mux := http.NewServeMux()
	mux.Handle("/swagger/", httpSwagger.WrapHandler)
	mux.HandleFunc("/api/login", handler.LoginHandler)
	mux.HandleFunc("/api/artikel/", handler.ArtikelHandler)

	// Urutan: Logging -> Auth -> Routing
	// protectedMux := BasicAuthMiddleware(mux)
	// loggedMux := LoggingMiddleware(protectedMux)

	protectedMux := JWTMiddleware(mux)
	loggedMux := LoggingMiddleware(protectedMux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      loggedMux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Dengarkan sinyal Ctrl+C (SIGINT) dan SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("Server REST API berjalan di http://localhost:%s...", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
		close(serverErr)
	}()

	select {
	case err := <-serverErr:
		if err != nil {
			log.Fatalf("Server berhenti karena error: %v", err)
		}
	case <-ctx.Done():
		log.Println("Sinyal shutdown diterima, menghentikan server...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("Graceful shutdown gagal, memaksa tutup server: %v", err)
			if closeErr := srv.Close(); closeErr != nil {
				log.Printf("Gagal menutup server paksa: %v", closeErr)
			}
		} else {
			log.Println("Server berhasil dihentikan dengan aman.")
		}
	}

	// fmt.Println("Server REST API berjalan di http://localhost:8080...")
	// err := http.ListenAndServe(":8080", loggedMux)
	// if err != nil {
	// 	fmt.Println("Gagal menyalakan server:", err)
	// }
}
