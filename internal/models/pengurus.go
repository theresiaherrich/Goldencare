package models

import (
	"time"

	"github.com/google/uuid"
)

type PengurusProfil struct {
	ID        uuid.UUID  `db:"id"         json:"id"`
	UserID    uuid.UUID  `db:"user_id"    json:"user_id"`
	KamarID   *uuid.UUID `db:"kamar_id"   json:"kamar_id"`
	Jabatan   string     `db:"jabatan"    json:"jabatan"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
}

type Shift struct {
	ID         uuid.UUID `db:"id"          json:"id"`
	UserID     uuid.UUID `db:"user_id"      json:"user_id"`
	NamaShift  string    `db:"nama_shift"   json:"nama_shift"`
	JamMulai   string    `db:"jam_mulai"    json:"jam_mulai"`
	JamSelesai string    `db:"jam_selesai"  json:"jam_selesai"`
	Hari       string    `db:"hari"         json:"hari"`
	CreatedAt  time.Time `db:"created_at"   json:"created_at"`
}

type CreatePengurusProfilRequest struct {
	KamarID *uuid.UUID `json:"kamar_id"`
	Jabatan string     `json:"jabatan"`
	Shifts  []struct {
		NamaShift  string `json:"nama_shift"  validate:"required"`
		JamMulai   string `json:"jam_mulai"   validate:"required"`
		JamSelesai string `json:"jam_selesai" validate:"required"`
		Hari       string `json:"hari"        validate:"required"`
	} `json:"shifts"`
}

type CatatanShift struct {
	ID         uuid.UUID `db:"id"            json:"id"`
	LansiaID   uuid.UUID `db:"lansia_id"     json:"lansia_id"`
	PengurusID uuid.UUID `db:"pengurus_id"   json:"pengurus_id"`
	IsiCatatan string    `db:"isi_catatan"   json:"isi_catatan"`
	Shift      string    `db:"shift"         json:"shift"`
	SuasanaHati  string `db:"suasana_hati"  json:"suasana_hati"`
	NafsuMakan   string `db:"nafsu_makan"   json:"nafsu_makan"`
	Aktivitas    string `db:"aktivitas"     json:"aktivitas"`
	StatusJurnal string `db:"status_jurnal" json:"status_jurnal"` 
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	NamaPengurus string `db:"nama_pengurus" json:"nama_pengurus,omitempty"`
	NamaLansia   string `db:"nama_lansia"   json:"nama_lansia,omitempty"`
}

type CreateCatatanShiftRequest struct {
	LansiaID     uuid.UUID `json:"lansia_id"   validate:"required"`
	IsiCatatan   string    `json:"isi_catatan" validate:"required"`
	Shift        string    `json:"shift"       validate:"required,oneof=Pagi Siang Malam"`
	SuasanaHati  string    `json:"suasana_hati" validate:"oneof=Tenang Gelisah Lesu"`
	NafsuMakan   string    `json:"nafsu_makan"  validate:"oneof=Baik Sedang Buruk"`
	Aktivitas    string    `json:"aktivitas"    validate:"oneof=Aktif Istirahat"`
	StatusJurnal string    `json:"status_jurnal" validate:"oneof=draf terkirim"`
}

var SuasanaHatiOptions = []string{"Tenang", "Gelisah", "Lesu"}
var NafsuMakanOptions = []string{"Baik", "Sedang", "Buruk"}
var AktivitasOptions = []string{"Aktif", "Istirahat"}
var ShiftOptions = []string{"Pagi", "Siang", "Malam"}

type JadwalCheckup struct {
	ID         uuid.UUID `db:"id"         json:"id"`
	LansiaID   uuid.UUID `db:"lansia_id"  json:"lansia_id"`
	Tanggal    time.Time `db:"tanggal"    json:"tanggal"`
	Keterangan string    `db:"keterangan" json:"keterangan"`
	Status     string    `db:"status"     json:"status"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

type ShiftHariIni struct {
	NamaShift   string `db:"nama_shift"   json:"nama_shift"`
	JamMulai    string `db:"jam_mulai"    json:"jam_mulai"`
	JamSelesai  string `db:"jam_selesai"  json:"jam_selesai"`
	JumlahStaf  int    `db:"jumlah_staf"  json:"jumlah_staf"`
}

type PengurusDashboard struct {
	TotalStaf      int            `json:"total_staf"`
	StafTanpaShift int            `json:"staf_tanpa_shift"`
	ShiftHariIni   []ShiftHariIni `json:"shift_hari_ini"`
	Hari           string         `json:"hari"`
}

type PengurusDetail struct {
	UserID    interface{} `db:"user_id"    json:"user_id"`
	Nama      string      `db:"nama"       json:"nama"`
	Email     string      `db:"email"      json:"email"`
	Jabatan   string      `db:"jabatan"    json:"jabatan"`
	NamaKamar string      `db:"nama_kamar" json:"nama_kamar"`
	NamaShift string      `db:"nama_shift" json:"nama_shift"`
	JamMulai  string      `db:"jam_mulai"  json:"jam_mulai"`
	JamSelesai string     `db:"jam_selesai" json:"jam_selesai"`
}

type PengurusDetailResponse struct {
	User      User           `json:"user"`
	Profil    PengurusProfil `json:"profil"`
	Shifts    []Shift        `json:"shifts"`
	NamaKamar string         `json:"nama_kamar"`
}

type KamarWithStaf struct {
	Kamar
	JumlahStaf   int `db:"jumlah_staf"   json:"jumlah_staf"`
	JumlahLansia int `db:"jumlah_lansia" json:"jumlah_lansia"`
}