package models

import (
	"time"

	"github.com/google/uuid"
)

type GaleriFisik struct {
	ID               uuid.UUID `db:"id"               json:"id"`
	LansiaID         uuid.UUID `db:"lansia_id"        json:"lansia_id"`
	PengurusID       uuid.UUID `db:"pengurus_id"      json:"pengurus_id"`
	FotoURL          string    `db:"foto_url"         json:"foto_url"`
	LokasiLuka       string    `db:"lokasi_luka"      json:"lokasi_luka"`
	Deskripsi        string    `db:"deskripsi"        json:"deskripsi"`
	AnalisisAI       string    `db:"analisis_ai"      json:"analisis_ai"`
	TingkatDarurat   string    `db:"tingkat_darurat"  json:"tingkat_darurat"`
	PrediksiPenyakit string    `db:"prediksi_penyakit" json:"prediksi_penyakit"`
	CreatedAt        time.Time `db:"created_at"       json:"created_at"`
}

type CreateGaleriRequest struct {
	LansiaID   uuid.UUID `json:"lansia_id"  validate:"required"`
	FotoURL    string    `json:"foto_url"`
	LokasiLuka string    `json:"lokasi_luka"`
	Deskripsi  string    `json:"deskripsi"`
}
