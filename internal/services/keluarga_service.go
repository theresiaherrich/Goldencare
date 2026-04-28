package service

import (
	"context"

	"github.com/theresiaherrich/Goldencare/internal/models"
	"github.com/theresiaherrich/Goldencare/internal/repository"
)

type KeluargaService interface {
	GetDashboard(ctx context.Context, lansiaID string) (*models.KeluargaDashboard, error)
	GetKeluarga(ctx context.Context, lansiaID string) (*models.KeluargaLengkap, error)
}

type keluargaService struct {
	repo repository.KeluargaRepository
}

func NewKeluargaService(repo repository.KeluargaRepository) KeluargaService {
	return &keluargaService{repo: repo}
}

func (s *keluargaService) GetDashboard(ctx context.Context, lansiaID string) (*models.KeluargaDashboard, error) {
	return s.repo.GetDashboard(ctx, lansiaID)
}

func (s *keluargaService) GetKeluarga(ctx context.Context, lansiaID string) (*models.KeluargaLengkap, error) {
	return s.repo.GetKeluarga(ctx, lansiaID)
}