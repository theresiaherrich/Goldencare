package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/theresiaherrich/Goldencare/internal/models"
	"github.com/theresiaherrich/Goldencare/internal/repository"
)

type PengurusService interface {
	GetDashboard(ctx context.Context, pantiID string) (*models.PengurusDashboard, error)
	GetAll(ctx context.Context, pantiID string, filterShift string) ([]models.PengurusDetail, error)
	GetByID(ctx context.Context, userID string, pantiID string) (*models.PengurusDetailResponse, error)
	SetProfil(ctx context.Context, req *models.CreatePengurusProfilRequest, targetUserID string) error
	GetShiftSaya(ctx context.Context, userID string) ([]models.Shift, error)
	GetKamar(ctx context.Context, pantiID string) ([]models.KamarWithStaf, error)
	CreateKamar(ctx context.Context, pantiID string, namaKamar string, kapasitas int) (*models.Kamar, error)
}

type pengurusService struct {
	repo repository.PengurusRepository
}

func NewPengurusService(repo repository.PengurusRepository) PengurusService {
	return &pengurusService{repo: repo}
}

func (s *pengurusService) GetDashboard(ctx context.Context, pantiID string) (*models.PengurusDashboard, error) {
	return s.repo.GetDashboard(ctx, pantiID)
}

func (s *pengurusService) GetAll(ctx context.Context, pantiID string, filterShift string) ([]models.PengurusDetail, error) {
	return s.repo.GetAll(ctx, pantiID, filterShift)
}

func (s *pengurusService) GetByID(ctx context.Context, userID string, pantiID string) (*models.PengurusDetailResponse, error) {
	return s.repo.GetByID(ctx, userID, pantiID)
}

func (s *pengurusService) SetProfil(ctx context.Context, req *models.CreatePengurusProfilRequest, targetUserID string) error {
	profil := &models.PengurusProfil{
		ID:        uuid.New(),
		UserID:    uuid.MustParse(targetUserID),
		KamarID:   req.KamarID,
		Jabatan:   req.Jabatan,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	var shifts []models.Shift
	for _, s := range req.Shifts {
		shifts = append(shifts, models.Shift{
			ID:         uuid.New(),
			UserID:     uuid.MustParse(targetUserID),
			NamaShift:  s.NamaShift,
			JamMulai:   s.JamMulai,
			JamSelesai: s.JamSelesai,
			Hari:       s.Hari,
			CreatedAt:  time.Now(),
		})
	}

	return s.repo.SetProfil(ctx, profil, shifts)
}

func (s *pengurusService) GetShiftSaya(ctx context.Context, userID string) ([]models.Shift, error) {
	return s.repo.GetShiftByUser(ctx, userID)
}

func (s *pengurusService) GetKamar(ctx context.Context, pantiID string) ([]models.KamarWithStaf, error) {
	return s.repo.GetKamar(ctx, pantiID)
}

func (s *pengurusService) CreateKamar(ctx context.Context, pantiID string, namaKamar string, kapasitas int) (*models.Kamar, error) {
	kamar := &models.Kamar{
		ID:        uuid.New(),
		PantiID:   uuid.MustParse(pantiID),
		NamaKamar: namaKamar,
		Kapasitas: kapasitas,
		CreatedAt: time.Now(),
	}
	if err := s.repo.CreateKamar(ctx, kamar); err != nil {
		return nil, err
	}
	return kamar, nil
}