package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-blog/repository"
)

type mockArtikelRepo struct {
	getAllFn  func(limit, offset int) ([]repository.Artikel, error)
	insertFn  func(judul, isi string) (int64, error)
	getByIDFn func(idStr string) (repository.Artikel, error)
}

func (m mockArtikelRepo) GetAllArtikel(limit, offset int) ([]repository.Artikel, error) {
	return m.getAllFn(limit, offset)
}

func (m mockArtikelRepo) InsertArtikel(judul, isi string) (int64, error) {
	return m.insertFn(judul, isi)
}

func (m mockArtikelRepo) GetArtikelByID(idStr string) (repository.Artikel, error) {
	return m.getByIDFn(idStr)
}

func TestArtikelHandler_PostSuccess_WithMock(t *testing.T) {
	origRepo := artikelRepo
	defer func() { artikelRepo = origRepo }()

	artikelRepo = mockArtikelRepo{
		insertFn: func(judul, isi string) (int64, error) {
			return 42, nil
		},
	}

	body := bytes.NewBufferString(`{"judul":"Belajar Go","isi":"Mock repository itu mantap"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/artikel/", body)
	rec := httptest.NewRecorder()

	ArtikelHandler(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	var resp map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("gagal parse JSON: %v", err)
	}

	if resp["status"] != true {
		t.Errorf("expected status=true, got %v", resp["status"])
	}
}
