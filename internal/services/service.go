package service

import (
	"github.com/theresiaherrich/Goldencare/internal/config"
	"github.com/theresiaherrich/Goldencare/internal/repository"
)

type Service interface {
	Auth() AuthService
	User() UserService
	Lansia() LansiaService
	Panti() PantiService
	Pengurus() PengurusService
	Kunjungan() KunjunganService
	Obat() ObatService
	Keluarga() KeluargaService
	Galeri() GaleriService
	CatatanShift() CatatanShiftService
	Pemeriksaan() PemeriksaanService
	Close() error
}

type serviceImpl struct {
	auth         AuthService
	user         UserService
	lansia       LansiaService
	panti        PantiService
	galeri       GaleriService
	catatanShift CatatanShiftService
	pemeriksaan  PemeriksaanService
	pengurus     PengurusService
	kunjungan    KunjunganService
	obat         ObatService
	keluarga     KeluargaService
}

func NewService(repo repository.Repository, cfg *config.Config) Service {
	return &serviceImpl{
		auth:         NewAuthService(repo, cfg),
		user:         NewUserService(repo),
		lansia:       NewLansiaService(repo),
		panti:        NewPantiService(repo),
		galeri:       NewGaleriService(repo.Galeri(), cfg),
		catatanShift: NewCatatanShiftService(repo.CatatanShift()),
		pemeriksaan:  NewPemeriksaanService(repo.Pemeriksaan()),
		pengurus:     NewPengurusService(repo.Pengurus()),
		kunjungan:    NewKunjunganService(repo.Kunjungan()),
		obat:         NewObatService(repo.Obat()),
		keluarga:     NewKeluargaService(repo.Keluarga()),
	}
}

func (s *serviceImpl) Auth() AuthService                 { return s.auth }
func (s *serviceImpl) User() UserService                 { return s.user }
func (s *serviceImpl) Lansia() LansiaService             { return s.lansia }
func (s *serviceImpl) Panti() PantiService               { return s.panti }
func (s *serviceImpl) Galeri() GaleriService             { return s.galeri }
func (s *serviceImpl) CatatanShift() CatatanShiftService { return s.catatanShift }
func (s *serviceImpl) Pemeriksaan() PemeriksaanService   { return s.pemeriksaan }
func (s *serviceImpl) Pengurus() PengurusService         { return s.pengurus }
func (s *serviceImpl) Kunjungan() KunjunganService       { return s.kunjungan }
func (s *serviceImpl) Obat() ObatService                 { return s.obat }
func (s *serviceImpl) Keluarga() KeluargaService         { return s.keluarga }

func (s *serviceImpl) Close() error { return nil }
