package models

import (
	"time"

	"github.com/google/uuid"
)

type Lansia struct {
	ID              uuid.UUID  `db:"id"               json:"id"`
	PantiID         uuid.UUID  `db:"panti_id"         json:"panti_id"`
	KamarID         *uuid.UUID `db:"kamar_id"         json:"kamar_id"`
	Nama            string     `db:"nama"             json:"nama"`
	NIK             string     `db:"nik"              json:"nik"`
	TanggalLahir    *time.Time `db:"tanggal_lahir"    json:"tanggal_lahir"`
	JenisKelamin    string     `db:"jenis_kelamin"    json:"jenis_kelamin"`
	AlamatAsal      string     `db:"alamat_asal"      json:"alamat_asal"`
	GolonganDarah   string     `db:"golongan_darah"   json:"golongan_darah"`
	RiwayatPenyakit string     `db:"riwayat_penyakit" json:"riwayat_penyakit"`
	Alergi          string     `db:"alergi"           json:"alergi"`
	FotoURL         string     `db:"foto_url"         json:"foto_url"`
	TanggalMasuk    time.Time  `db:"tanggal_masuk"    json:"tanggal_masuk"`
	Status          string     `db:"status"           json:"status"`
	CreatedAt       time.Time  `db:"created_at"       json:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at"       json:"updated_at"`
}

type CreateLansiaRequest struct {
	KamarID         *uuid.UUID `json:"kamar_id"`
	Nama            string     `json:"nama"             validate:"required"`
	NIK             string     `json:"nik"`
	TanggalLahir    string     `json:"tanggal_lahir"` 
	JenisKelamin    string     `json:"jenis_kelamin"    validate:"oneof=L P"`
	AlamatAsal      string     `json:"alamat_asal"`
	GolonganDarah   string     `json:"golongan_darah"`
	RiwayatPenyakit string     `json:"riwayat_penyakit"`
	Alergi          string     `json:"alergi"`
}

type AktivitasLansia struct {
	ID        uuid.UUID `db:"id"         json:"id"`
	LansiaID  uuid.UUID `db:"lansia_id"  json:"lansia_id"`
	Nama      string    `db:"nama"       json:"nama"`
	Deskripsi string    `db:"deskripsi"  json:"deskripsi"`
	Jam       string    `db:"jam"        json:"jam"`
	Hari      string    `db:"hari"       json:"hari"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
