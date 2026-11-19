package handler

import (
	"encoding/json"
	"fmt"
	"go-blog/repository"
	"net/http"
	"strings"
	// Impor repository yang baru dibuat
)

// Helper Error & Sukses tetap di sini...
func KirimError(w http.ResponseWriter, pesan string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]any{
		"status": false,
		"code":   statusCode,
		"error":  pesan,
	})
}

func KirimSukses(w http.ResponseWriter, pesan string, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]any{
		"status": true,
		"code":   statusCode,
		"pesan":  pesan,
		"data":   data,
	})
}

func ArtikelHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/api/artikel")
	if path != "" && path != "/" {
		handleArtikelDetail(w, r, path)
		return
	}

	// 1. METHOD GET (Menggunakan Repository)
	if r.Method == http.MethodGet {
		limit := 10
		page := 1
		offset := 0

		if qLimit := r.URL.Query().Get("limit"); qLimit != "" {
			fmt.Sscanf(qLimit, "%d", &limit)
		}
		if qPage := r.URL.Query().Get("page"); qPage != "" {
			var p int
			fmt.Sscanf(qPage, "%d", &p)
			if p > 0 {
				page = p
				offset = (page - 1) * limit
			}
		}

		// Panggil repository
		daftarArtikel, err := artikelRepo.GetAllArtikel(limit, offset)
		if err != nil {
			KirimError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		KirimSukses(w, "Berhasil mengambil data artikel", map[string]any{
			"page":  page,
			"limit": limit,
			"data":  daftarArtikel,
		}, http.StatusOK)
		return
	}

	// 2. METHOD POST (Menggunakan Repository)
	if r.Method == http.MethodPost {
		var artikelBaru repository.Artikel
		err := json.NewDecoder(r.Body).Decode(&artikelBaru)
		if err != nil {
			KirimError(w, "Format JSON tidak valid", http.StatusBadRequest)
			return
		}

		if artikelBaru.Judul == "" || artikelBaru.Isi == "" {
			KirimError(w, "Judul dan isi artikel tidak boleh kosong!", http.StatusBadRequest)
			return
		}

		// Panggil repository untuk insert
		lastID, err := artikelRepo.InsertArtikel(artikelBaru.Judul, artikelBaru.Isi)
		if err != nil {
			KirimError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		artikelBaru.ID = int(lastID)
		KirimSukses(w, "Artikel berhasil disimpan ke MariaDB! 🎉", artikelBaru, http.StatusCreated)
		return
	}
}

func handleArtikelDetail(w http.ResponseWriter, r *http.Request, path string) {
	idStr := strings.TrimPrefix(path, "/")

	if r.Method == http.MethodGet {
		a, err := artikelRepo.GetArtikelByID(idStr)
		if err != nil {
			KirimError(w, "Artikel tidak ditemukan", http.StatusNotFound)
			return
		}
		KirimSukses(w, "Detail artikel", a, http.StatusOK)
		return
	}
}
