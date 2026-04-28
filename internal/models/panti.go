package models

import (
	"time"

	"github.com/google/uuid"
)

type Panti struct {
	ID           uuid.UUID `db:"id" json:"id"`
	Nama         string    `db:"nama" json:"nama"`
	Alamat       *string   `db:"alamat" json:"alamat"`
	Telepon      *string   `db:"telepon" json:"telepon"`
	KodeUndangan string    `db:"kode_undangan" json:"kode_undangan"`
	PengelolaID  uuid.UUID `db:"pengelola_id" json:"pengelola_id"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type CreatePantiRequest struct {
	Nama    string `json:"nama"    validate:"required"`
	Alamat  string `json:"alamat"`
	Telepon string `json:"telepon"`
}

type Kamar struct {
	ID        uuid.UUID `db:"id"         json:"id"`
	PantiID   uuid.UUID `db:"panti_id"   json:"panti_id"`
	NamaKamar string    `db:"nama_kamar" json:"nama_kamar"`
	Kapasitas int       `db:"kapasitas"  json:"kapasitas"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}