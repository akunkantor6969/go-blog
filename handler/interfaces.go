package handler

import "go-blog/repository"

type ArtikelRepository interface {
	GetAllArtikel(limit, offset int) ([]repository.Artikel, error)
	InsertArtikel(judul, isi string) (int64, error)
	GetArtikelByID(idStr string) (repository.Artikel, error)
}

type defaultArtikelRepo struct{}

func (defaultArtikelRepo) GetAllArtikel(limit, offset int) ([]repository.Artikel, error) {
	return repository.GetAllArtikel(limit, offset)
}

func (defaultArtikelRepo) InsertArtikel(judul, isi string) (int64, error) {
	return repository.InsertArtikel(judul, isi)
}

func (defaultArtikelRepo) GetArtikelByID(idStr string) (repository.Artikel, error) {
	return repository.GetArtikelByID(idStr)
}

var artikelRepo ArtikelRepository = defaultArtikelRepo{}
