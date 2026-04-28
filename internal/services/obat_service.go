package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/theresiaherrich/Goldencare/internal/models"
	"github.com/theresiaherrich/Goldencare/internal/repository"
)

type ObatService interface {
	Create(ctx context.Context, req *models.CreateObatRequest) (*models.ObatDetail, error)
	GetByLansia(ctx context.Context, lansiaID string) ([]models.ObatDetail, error)
	GetRiwayatPemberian(ctx context.Context, lansiaID string) ([]models.RiwayatPemberianItem, error)
	GetJadwalHariIni(ctx context.Context, pengurusID string) (*models.JadwalHariIniResponse, error)
	Checklist(ctx context.Context, jadwalObatID string, pengurusID string, catatan string) (*models.LogPemberianObat, error)
	DeleteObat(ctx context.Context, obatID string) error
}

type obatService struct {
	repo repository.ObatRepository
}

func NewObatService(repo repository.ObatRepository) ObatService {
	return &obatService{repo: repo}
}

func (s *obatService) Create(ctx context.Context, req *models.CreateObatRequest) (*models.ObatDetail, error) {
	if len(req.Jadwals) == 0 {
		return nil, errors.New("minimal 1 jadwal pemberian obat harus diisi")
	}

	now := time.Now()
	obat := &models.Obat{
		ID:            uuid.New(),
		LansiaID:      req.LansiaID,
		NamaObat:      req.NamaObat,
		Dosis:         req.Dosis,
		CaraPemberian: req.CaraPemberian,
		Keterangan:    req.Keterangan,
		IsAktif:       true,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	var jadwals []models.JadwalObat
	for _, j := range req.Jadwals {
		jadwals = append(jadwals, models.JadwalObat{
			ID:        uuid.New(),
			ObatID:    obat.ID,
			Jam:       j.Jam,
			Shift:     j.Shift,
			CreatedAt: now,
		})
	}

	if err := s.repo.CreateWithJadwal(ctx, obat, jadwals); err != nil {
		return nil, err
	}

	return &models.ObatDetail{Obat: *obat, Jadwals: jadwals}, nil
}

func (s *obatService) GetByLansia(ctx context.Context, lansiaID string) ([]models.ObatDetail, error) {
	return s.repo.GetByLansia(ctx, lansiaID)
}

func (s *obatService) GetRiwayatPemberian(ctx context.Context, lansiaID string) ([]models.RiwayatPemberianItem, error) {
	return s.repo.GetRiwayatPemberian(ctx, lansiaID)
}

func (s *obatService) GetJadwalHariIni(ctx context.Context, pengurusID string) (*models.JadwalHariIniResponse, error) {
	items, err := s.repo.GetJadwalHariIni(ctx, pengurusID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	grouped := map[string][]models.JadwalHariIniItem{
		"Pagi":  {},
		"Siang": {},
		"Sore":  {},
	}
	for _, item := range items {
		item.StatusPemberian = hitungStatusPemberian(item.Jam, item.DiberikanPada, now)
		grouped[item.Shift] = append(grouped[item.Shift], item)
	}

	return &models.JadwalHariIniResponse{
		Tanggal: now.Format("2006-01-02"),
		Hari:    hariIndonesia(now.Weekday()),
		Jadwal:  grouped,
	}, nil
}

func (s *obatService) Checklist(ctx context.Context, jadwalObatID string, pengurusID string, catatan string) (*models.LogPemberianObat, error) {
	sudahAda, err := s.repo.CekSudahDiberikan(ctx, jadwalObatID)
	if err != nil {
		return nil, err
	}
	if sudahAda {
		return nil, errors.New("obat ini sudah ditandai diberikan hari ini")
	}

	log := &models.LogPemberianObat{
		ID:            uuid.New(),
		JadwalObatID:  uuid.MustParse(jadwalObatID),
		PengurusID:    uuid.MustParse(pengurusID),
		DiberikanPada: time.Now(),
		Catatan:       catatan,
	}

	if err := s.repo.CreateLog(ctx, log); err != nil {
		return nil, err
	}
	return log, nil
}

func (s *obatService) DeleteObat(ctx context.Context, obatID string) error {
	return s.repo.SoftDelete(ctx, obatID)
}