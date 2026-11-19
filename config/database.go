package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var DB *sql.DB

func InitDatabase() {
	// Muat file .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: File .env tidak ditemukan, menggunakan variabel sistem default.")
	}

	// Ambil nilai dari environment variables
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Susun string koneksi DSN secara dinamis
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Gagal membuka koneksi database: ", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Gagal terhubung ke database MariaDB: ", err)
	}

	DB = db
	fmt.Println("Berhasil terhubung ke MariaDB menggunakan konfigurasi .env! 🐬")
}
