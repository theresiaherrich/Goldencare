package models

import (
	"time"

	"github.com/google/uuid"
)

type PemeriksaanKesehatan struct {
	ID            uuid.UUID `db:"id"             json:"id"`
	LansiaID      uuid.UUID `db:"lansia_id"      json:"lansia_id"`
	PengurusID    uuid.UUID `db:"pengurus_id"    json:"pengurus_id"`
	TekananDarah  string    `db:"tekanan_darah"  json:"tekanan_darah"`
	DetakJantung  int       `db:"detak_jantung"  json:"detak_jantung"`
	GulaDarah     float64   `db:"gula_darah"     json:"gula_darah"`
	SuhuTubuh     float64   `db:"suhu_tubuh"     json:"suhu_tubuh"`
	BeratBadan    float64   `db:"berat_badan"    json:"berat_badan"`
	Keluhan       string    `db:"keluhan"        json:"keluhan"`
	Status        string    `db:"status"         json:"status"`
	Rekomendasi   string    `db:"rekomendasi"    json:"rekomendasi"`
	StatusDarurat string    `db:"status_darurat" json:"status_darurat"`
	CreatedAt     time.Time `db:"created_at"     json:"created_at"`
}

type CreatePemeriksaanRequest struct {
	LansiaID     uuid.UUID `json:"lansia_id"     validate:"required"`
	TekananDarah string    `json:"tekanan_darah"`
	DetakJantung int       `json:"detak_jantung"`
	GulaDarah    float64   `json:"gula_darah"`
	SuhuTubuh    float64   `json:"suhu_tubuh"`
	BeratBadan   float64   `json:"berat_badan"`
	Keluhan      string    `json:"keluhan"`
}
