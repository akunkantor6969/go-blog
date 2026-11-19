# go-blog

REST API blog sederhana berbasis Go (`net/http`) dengan arsitektur modular, autentikasi JWT, *unit testing* menggunakan *mocking* repository, dan *graceful shutdown* dengan koneksi MariaDB.

## Fitur Utama

- **Arsitektur Modular**: Pemisahan tanggung jawab ke dalam package `config`, `handler`, dan `repository`.
- **Autentikasi JWT**: Endpoint login (`POST /api/login`) menggunakan `github.com/golang-jwt/jwt/v5` dan proteksi *request* non-`GET` lewat `JWTMiddleware`.
- **Unit Testing & Mocking**: Pengujian handler, validasi, login, dan *mock* `ArtikelRepository` menggunakan `testing` dan `net/http/httptest`.
- **Graceful Shutdown**: Menangani sinyal berhenti (`SIGINT` / `SIGTERM`) agar server dan koneksi database ditutup dengan rapi.

## Struktur Folder

```text
go-blog/
├── config/          # Inisialisasi konfigurasi dan koneksi MariaDB
├── handler/         # HTTP handler, response helper, interface repo, dan unit test
├── repository/      # Akses data dan query ke MariaDB
├── .env.example     # Contoh variabel environment
├── .gitignore       # Daftar file/folder yang diabaikan Git
├── go.mod           # Definisi module dan dependensi Go
├── go.sum           # Checksum dependensi
└── main.go          # Entry point, middleware, routing, dan graceful shutdown