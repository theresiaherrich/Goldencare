package models

import (
	"time"

	"github.com/google/uuid"
)

type Obat struct {
	ID            uuid.UUID `db:"id"            json:"id"`
	LansiaID      uuid.UUID `db:"lansia_id"     json:"lansia_id"`
	NamaObat      string    `db:"nama_obat"     json:"nama_obat"`
	Dosis         string    `db:"dosis"         json:"dosis"`
	CaraPemberian string    `db:"cara_pemberian" json:"cara_pemberian"` 
	Keterangan    string    `db:"keterangan"    json:"keterangan"`
	IsAktif       bool      `db:"is_aktif"      json:"is_aktif"`
	CreatedAt     time.Time `db:"created_at"    json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"    json:"updated_at"`
}

type JadwalObat struct {
	ID        uuid.UUID `db:"id"         json:"id"`
	ObatID    uuid.UUID `db:"obat_id"    json:"obat_id"`
	Jam       string    `db:"jam"        json:"jam"`   
	Shift     string    `db:"shift"      json:"shift"` 
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type LogPemberianObat struct {
	ID            uuid.UUID `db:"id"             json:"id"`
	JadwalObatID  uuid.UUID `db:"jadwal_obat_id" json:"jadwal_obat_id"`
	PengurusID    uuid.UUID `db:"pengurus_id"    json:"pengurus_id"`
	DiberikanPada time.Time `db:"diberikan_pada" json:"diberikan_pada"`
	Catatan       string    `db:"catatan"        json:"catatan"`
}

type CreateObatRequest struct {
	LansiaID      uuid.UUID `json:"lansia_id"      validate:"required"`
	NamaObat      string    `json:"nama_obat"      validate:"required"`
	Dosis         string    `json:"dosis"`
	CaraPemberian string    `json:"cara_pemberian" validate:"required,oneof=Oral Inhaler Injeksi Tetes Topikal"`
	Keterangan    string    `json:"keterangan"`
	Jadwals       []struct {
		Jam   string `json:"jam"   validate:"required"` 
		Shift string `json:"shift" validate:"required,oneof=Pagi Siang Sore"`
	} `json:"jadwals" validate:"required,min=1"`
}

var CaraPemberianOptions = []string{"Oral", "Inhaler", "Injeksi", "Tetes", "Topikal"}
var ShiftObatOptions = []string{"Pagi", "Siang", "Sore"}


type JadwalHariIniItem struct {
	JadwalObatID    uuid.UUID  `db:"jadwal_obat_id"  json:"jadwal_obat_id"`
	ObatID          uuid.UUID  `db:"obat_id"         json:"obat_id"`
	LansiaID        uuid.UUID  `db:"lansia_id"       json:"lansia_id"`
	NamaLansia      string     `db:"nama_lansia"     json:"nama_lansia"`
	NomorKamar      string     `db:"nomor_kamar"     json:"nomor_kamar"`
	FotoLansia      string     `db:"foto_lansia"     json:"foto_lansia"`
	NamaObat        string     `db:"nama_obat"       json:"nama_obat"`
	Dosis           string     `db:"dosis"           json:"dosis"`
	CaraPemberian   string     `db:"cara_pemberian"  json:"cara_pemberian"`
	Jam             string     `db:"jam"             json:"jam"`
	Shift           string     `db:"shift"           json:"shift"`
	StatusPemberian string     `db:"-"               json:"status_pemberian"` 
	DiberikanPada   *time.Time `db:"diberikan_pada"  json:"diberikan_pada"`
	NamaPengurus    string     `db:"nama_pengurus"   json:"nama_pengurus"`
}

type JadwalHariIniResponse struct {
	Tanggal string                         `json:"tanggal"`
	Hari    string                         `json:"hari"`
	Jadwal  map[string][]JadwalHariIniItem `json:"jadwal"`
}

type ObatDetail struct {
	Obat    Obat         `json:"obat"`
	Jadwals []JadwalObat `json:"jadwals"`
}


type RiwayatPemberianItem struct {
	ID            interface{} `db:"id"             json:"id"`
	JadwalObatID  interface{} `db:"jadwal_obat_id" json:"jadwal_obat_id"`
	PengurusID    interface{} `db:"pengurus_id"    json:"pengurus_id"`
	DiberikanPada time.Time   `db:"diberikan_pada" json:"diberikan_pada"`
	Catatan       string      `db:"catatan"        json:"catatan"`
	NamaObat      string `db:"nama_obat"      json:"nama_obat"`
	Dosis         string `db:"dosis"          json:"dosis"`
	CaraPemberian string `db:"cara_pemberian" json:"cara_pemberian"`
	Jam           string `db:"jam"            json:"jam"`
	Shift         string `db:"shift"          json:"shift"`
	NamaPengurus string `db:"nama_pengurus" json:"nama_pengurus"`
}