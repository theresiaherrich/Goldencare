package service

import (
	"context"
	"fmt"

	"github.com/theresiaherrich/Goldencare/internal/models"
	"github.com/theresiaherrich/Goldencare/internal/repository"
)

type LansiaService interface {
	Create(ctx context.Context, lansia *models.Lansia) error
	GetByID(ctx context.Context, id string) (*models.Lansia, error)
	GetByPantiID(ctx context.Context, pantiID string, filters map[string]interface{}) ([]models.Lansia, error)
	GetDashboard(ctx context.Context, pantiID string) (map[string]interface{}, error)
	Update(ctx context.Context, lansia *models.Lansia) error
}

type lansiaService struct {
	repo repository.Repository
}

func NewLansiaService(repo repository.Repository) LansiaService {
	return &lansiaService{repo: repo}
}

func (s *lansiaService) Create(ctx context.Context, lansia *models.Lansia) error {
	return s.repo.Lansia().Create(ctx, lansia)
}

func (s *lansiaService) GetByID(ctx context.Context, id string) (*models.Lansia, error) {
	return s.repo.Lansia().GetByID(ctx, id)
}

func (s *lansiaService) GetByPantiID(ctx context.Context, pantiID string, filters map[string]interface{}) ([]models.Lansia, error) {
	return s.repo.Lansia().GetByPantiID(ctx, pantiID, filters)
}

func (s *lansiaService) GetDashboard(ctx context.Context, pantiID string) (map[string]interface{}, error) {
	totalLansia, err := s.repo.Lansia().GetCountByPantiID(ctx, pantiID)
	if err != nil {
		return nil, fmt.Errorf("gagal menghitung total lansia: %w", err)
	}
	return map[string]interface{}{
		"total_lansia":         totalLansia,
		"perlu_perhatian":      0,
		"pemeriksaan_hari_ini": 0,
	}, nil
}

func (s *lansiaService) Update(ctx context.Context, lansia *models.Lansia) error {
	return s.repo.Lansia().Update(ctx, lansia)
}