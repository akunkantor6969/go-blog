package repository

import (
	"go-blog/config"
)

type Artikel struct {
	ID    int    `json:"id"`
	Judul string `json:"judul"`
	Isi   string `json:"isi"`
}

// Mengambil semua artikel dengan pagination
func GetAllArtikel(limit, offset int) ([]Artikel, error) {
	query := "SELECT id, judul, isi FROM artikel LIMIT ? OFFSET ?"
	rows, err := config.DB.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var daftarArtikel []Artikel
	for rows.Next() {
		var a Artikel
		err := rows.Scan(&a.ID, &a.Judul, &a.Isi)
		if err != nil {
			return nil, err
		}
		daftarArtikel = append(daftarArtikel, a)
	}

	if daftarArtikel == nil {
		daftarArtikel = []Artikel{}
	}

	return daftarArtikel, nil
}

// Menyimpan artikel baru ke database
func InsertArtikel(judul, isi string) (int64, error) {
	query := "INSERT INTO artikel (judul, isi) VALUES (?, ?)"
	result, err := config.DB.Exec(query, judul, isi)
	if err != nil {
		return 0, err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastID, nil
}

// Mengambil artikel berdasarkan ID
func GetArtikelByID(idStr string) (Artikel, error) {
	var a Artikel
	query := "SELECT id, judul, isi FROM artikel WHERE id = ?"
	err := config.DB.QueryRow(query, idStr).Scan(&a.ID, &a.Judul, &a.Isi)
	if err != nil {
		return a, err
	}
	return a, nil
}
