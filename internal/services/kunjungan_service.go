package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/theresiaherrich/Goldencare/internal/models"
	"github.com/theresiaherrich/Goldencare/internal/repository"
)

type KunjunganService interface {
	GetByLansia(ctx context.Context, lansiaID string) ([]models.KunjunganKeluarga, error)
	GetTerbaru(ctx context.Context, lansiaID string) ([]models.KunjunganKeluarga, error)
	Create(ctx context.Context, req *models.CreateKunjunganRequest, pengurusID string) (*models.KunjunganKeluarga, error)
}

type kunjunganService struct {
	repo repository.KunjunganRepository
}

func NewKunjunganService(repo repository.KunjunganRepository) KunjunganService {
	return &kunjunganService{repo: repo}
}

func (s *kunjunganService) GetByLansia(ctx context.Context, lansiaID string) ([]models.KunjunganKeluarga, error) {
	return s.repo.GetByLansia(ctx, lansiaID)
}

func (s *kunjunganService) GetTerbaru(ctx context.Context, lansiaID string) ([]models.KunjunganKeluarga, error) {
	return s.repo.GetTerbaru(ctx, lansiaID, 3)
}

func (s *kunjunganService) Create(ctx context.Context, req *models.CreateKunjunganRequest, pengurusID string) (*models.KunjunganKeluarga, error) {
	if req.LansiaID == uuid.Nil {
		return nil, errors.New("lansia_id wajib diisi")
	}
	if req.NamaKeluarga == "" {
		return nil, errors.New("nama_keluarga wajib diisi")
	}

	tanggalKunjungan := time.Now()
	if req.TanggalKunjungan != "" {
		parsed, err := time.Parse(time.RFC3339, req.TanggalKunjungan)
		if err != nil {
			return nil, errors.New("format tanggal kunjungan tidak valid, gunakan RFC3339")
		}
		tanggalKunjungan = parsed
	}

	kunjungan := &models.KunjunganKeluarga{
		ID:               uuid.New(),
		LansiaID:         req.LansiaID,
		PengurusID:       uuid.MustParse(pengurusID),
		NamaKeluarga:     req.NamaKeluarga,
		HubunganKeluarga: req.HubunganKeluarga,
		TanggalKunjungan: tanggalKunjungan,
		DurasiMenit:      req.DurasiMenit,
		FotoURL:          req.FotoURL,
		Catatan:          req.Catatan,
		ResponLansia:     req.ResponLansia,
		CreatedAt:        time.Now(),
	}

	if err := s.repo.Create(ctx, kunjungan); err != nil {
		return nil, err
	}
	return kunjungan, nil
}