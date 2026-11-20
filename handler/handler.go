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

	switch r.Method {
	case http.MethodGet:
		GetArtikelHandler(w, r)
	case http.MethodPost:
		CreateArtikelHandler(w, r)
	default:
		KirimError(w, "Method tidak diizinkan", http.StatusMethodNotAllowed)
	}
}

// GetArtikelHandler godoc
// @Summary      Ambil daftar atau detail artikel
// @Description  Mengambil daftar artikel dengan pagination (?page=1&limit=10) atau detail artikel berdasarkan ID
// @Tags         Artikel
// @Produce      json
// @Param        page   query     int  false  "Nomor halaman"
// @Param        limit  query     int  false  "Jumlah data per halaman"
// @Success      200    {object}  map[string]interface{}
// @Failure      400    {object}  map[string]interface{}
// @Router       /api/artikel/ [get]
func GetArtikelHandler(w http.ResponseWriter, r *http.Request) {
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
}

// CreateArtikelHandler godoc
// @Summary      Buat artikel baru
// @Description  Menambahkan artikel baru ke database (membutuhkan token JWT)
// @Tags         Artikel
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        artikel  body      repository.Artikel  true  "Data Artikel Baru"
// @Success      201      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      401      {object}  map[string]interface{}
// @Router       /api/artikel/ [post]
func CreateArtikelHandler(w http.ResponseWriter, r *http.Request) {
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

// func ArtikelHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	path := strings.TrimPrefix(r.URL.Path, "/api/artikel")
// 	if path != "" && path != "/" {
// 		handleArtikelDetail(w, r, path)
// 		return
// 	}

// 	// 1. METHOD GET (Menggunakan Repository)
// 	if r.Method == http.MethodGet {
// 		limit := 10
// 		page := 1
// 		offset := 0

// 		if qLimit := r.URL.Query().Get("limit"); qLimit != "" {
// 			fmt.Sscanf(qLimit, "%d", &limit)
// 		}
// 		if qPage := r.URL.Query().Get("page"); qPage != "" {
// 			var p int
// 			fmt.Sscanf(qPage, "%d", &p)
// 			if p > 0 {
// 				page = p
// 				offset = (page - 1) * limit
// 			}
// 		}

// 		// Panggil repository
// 		daftarArtikel, err := artikelRepo.GetAllArtikel(limit, offset)
// 		if err != nil {
// 			KirimError(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		KirimSukses(w, "Berhasil mengambil data artikel", map[string]any{
// 			"page":  page,
// 			"limit": limit,
// 			"data":  daftarArtikel,
// 		}, http.StatusOK)
// 		return
// 	}

// 	// 2. METHOD POST (Menggunakan Repository)
// 	if r.Method == http.MethodPost {
// 		var artikelBaru repository.Artikel
// 		err := json.NewDecoder(r.Body).Decode(&artikelBaru)
// 		if err != nil {
// 			KirimError(w, "Format JSON tidak valid", http.StatusBadRequest)
// 			return
// 		}

// 		if artikelBaru.Judul == "" || artikelBaru.Isi == "" {
// 			KirimError(w, "Judul dan isi artikel tidak boleh kosong!", http.StatusBadRequest)
// 			return
// 		}

// 		// Panggil repository untuk insert
// 		lastID, err := artikelRepo.InsertArtikel(artikelBaru.Judul, artikelBaru.Isi)
// 		if err != nil {
// 			KirimError(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		artikelBaru.ID = int(lastID)
// 		KirimSukses(w, "Artikel berhasil disimpan ke MariaDB! 🎉", artikelBaru, http.StatusCreated)
// 		return
// 	}
// }

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
