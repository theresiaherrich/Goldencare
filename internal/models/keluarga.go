package models

import (
	"time"

	"github.com/google/uuid"
)

type KunjunganKeluarga struct {
	ID               uuid.UUID `db:"id"                json:"id"`
	LansiaID         uuid.UUID `db:"lansia_id"         json:"lansia_id"`
	PengurusID       uuid.UUID `db:"pengurus_id"       json:"pengurus_id"`
	NamaKeluarga     string    `db:"nama_keluarga"     json:"nama_keluarga"`
	HubunganKeluarga string    `db:"hubungan_keluarga" json:"hubungan_keluarga"`
	TanggalKunjungan time.Time `db:"tanggal_kunjungan" json:"tanggal_kunjungan"`
	DurasiMenit      int       `db:"durasi_menit"      json:"durasi_menit"`
	FotoURL          string    `db:"foto_url"          json:"foto_url"`
	Catatan          string    `db:"catatan"           json:"catatan"`
	ResponLansia     string    `db:"respon_lansia"     json:"respon_lansia"`
	CreatedAt        time.Time `db:"created_at"        json:"created_at"`
	NamaPengurus     string    `db:"nama_pengurus" json:"nama_pengurus,omitempty"`
}

type CreateKunjunganRequest struct {
	LansiaID         uuid.UUID `json:"lansia_id"          validate:"required"`
	NamaKeluarga     string    `json:"nama_keluarga"      validate:"required"`
	HubunganKeluarga string    `json:"hubungan_keluarga"`
	TanggalKunjungan string    `json:"tanggal_kunjungan"`
	DurasiMenit      int       `json:"durasi_menit"`
	FotoURL          string    `json:"foto_url"`
	Catatan          string    `json:"catatan"`
	ResponLansia     string    `json:"respon_lansia"`
}

type CatatanShiftPublik struct {
	ID           interface{} `db:"id"            json:"id"`
	NarasiShift  string      `db:"narasi_shift"  json:"narasi_shift"`
	SuasanaHati  string      `db:"suasana_hati"  json:"suasana_hati"`
	NafsuMakan   string      `db:"nafsu_makan"   json:"nafsu_makan"`
	Aktivitas    string      `db:"aktivitas"     json:"aktivitas"`
	NamaShift    string      `db:"nama_shift"    json:"nama_shift"`
	NamaPengurus string      `db:"nama_pengurus" json:"nama_pengurus"`
	CreatedAt    time.Time   `db:"created_at"    json:"created_at"`
}

type VitalPublik struct {
	TekananDarah string    `db:"tekanan_darah" json:"tekanan_darah"`
	FaseLabel    string    `db:"fase_label"    json:"fase_label"`
	FaseWarna    string    `db:"fase_warna"    json:"fase_warna"`
	SuhuTubuh    float64   `db:"suhu_tubuh"    json:"suhu_tubuh"`
	GulaDarah    float64   `db:"gula_darah"    json:"gula_darah"`
	CreatedAt    time.Time `db:"created_at"    json:"created_at"`
}

type ObatStatus struct {
	NamaObat       string `db:"nama_obat"        json:"nama_obat"`
	Dosis          string `db:"dosis"            json:"dosis"`
	SudahDiberikan bool   `db:"sudah_diberikan"  json:"sudah_diberikan"`
	WaktuDiberikan string `db:"waktu_diberikan"  json:"waktu_diberikan"`
}

type KunjunganRingkas struct {
	NamaKeluarga     string    `db:"nama_keluarga"     json:"nama_keluarga"`
	HubunganKeluarga string    `db:"hubungan_keluarga" json:"hubungan_keluarga"`
	TanggalKunjungan time.Time `db:"tanggal_kunjungan" json:"tanggal_kunjungan"`
	DurasiMenit      int       `db:"durasi_menit"      json:"durasi_menit"`
	Catatan          string    `db:"catatan"           json:"catatan"`
}

type KondisiFisikRingkas struct {
	JenisKondisi  string    `db:"jenis_kondisi" json:"jenis_kondisi"`
	LokasiTubuh   string    `db:"lokasi_tubuh"  json:"lokasi_tubuh"`
	RisikoLabel   string    `db:"risiko_label"  json:"risiko_label"`
	TingkatRisiko string    `db:"tingkat_risiko" json:"tingkat_risiko"`
	FotoURL       string    `db:"foto_url"      json:"foto_url"`
	CreatedAt     time.Time `db:"created_at"    json:"created_at"`
}

type ObatHariIniSummary struct {
	Total          int `json:"total"`
	SudahDiberikan int `json:"sudah_diberikan"`
	BelumDiberikan int `json:"belum_diberikan"`
}

type KeluargaDashboard struct {
	Lansia              Lansia               `json:"lansia"`
	CatatanHariIni      []CatatanShiftPublik `json:"catatan_hari_ini"`
	VitalTerbaru        VitalPublik          `json:"vital_terbaru"`
	ObatPagi            []ObatStatus         `json:"obat_pagi"`
	RiwayatKunjungan    []KunjunganRingkas   `json:"riwayat_kunjungan"`
	KondisiFisikTerbaru KondisiFisikRingkas  `json:"kondisi_fisik_terbaru"`
}

type KeluargaLengkap struct {
	Lansia              Lansia                 `json:"lansia"`
	CatatanShift        []CatatanShift         `json:"catatan_shift"`
	RiwayatPemeriksaan  []PemeriksaanKesehatan `json:"riwayat_pemeriksaan"`
	ObatAktif           []Obat                 `json:"obat_aktif"`
	ObatHariIni         ObatHariIniSummary     `json:"obat_hari_ini"`
	RiwayatKondisiFisik []GaleriFisik          `json:"riwayat_kondisi_fisik"`
	RiwayatKunjungan    []KunjunganKeluarga    `json:"riwayat_kunjungan"`
}
