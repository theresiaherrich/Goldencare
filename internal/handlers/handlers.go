package handlers

import "github.com/theresiaherrich/Goldencare/internal/services"

type Handlers struct {
	Auth         *AuthHandler
	Panti        *PantiHandler
	Lansia       *LansiaHandler
	Pengurus     *PengurusHandler
	Kunjungan    *KunjunganHandler
	Obat         *ObatHandler
	Keluarga     *KeluargaHandler
	CatatanShift *CatatanShiftHandler
	Pemeriksaan  *PemeriksaanHandler
	Galeri       *GaleriHandler
}

func NewHandlers(svc service.Service) *Handlers {
	return &Handlers{
		Auth:         NewAuthHandler(svc.Auth()),
		Panti:        NewPantiHandler(svc.Panti()),
		Lansia:       NewLansiaHandler(svc.Lansia()),
		Pengurus:     NewPengurusHandler(svc.Pengurus()),
		Kunjungan:    NewKunjunganHandler(svc.Kunjungan()),
		Obat:         NewObatHandler(svc.Obat()),
		Keluarga:     NewKeluargaHandler(svc.Keluarga()),
		CatatanShift: NewCatatanShiftHandler(svc.CatatanShift()),
		Pemeriksaan:  NewPemeriksaanHandler(svc.Pemeriksaan()),
		Galeri:       NewGaleriHandler(svc.Galeri()),
	}
}
