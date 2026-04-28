package service

import (
	"context"
	"fmt"
	"time"

	"github.com/theresiaherrich/Goldencare/internal/models"
	"github.com/theresiaherrich/Goldencare/internal/repository"
)

type PantiService interface {
	GetByID(ctx context.Context, id string) (*models.Panti, error)
	GetAll(ctx context.Context) ([]models.Panti, error)
	Create(ctx context.Context, panti *models.Panti) error
	Update(ctx context.Context, panti *models.Panti) error
	UpdatePanti(ctx context.Context, userID string, pantiID string, update map[string]interface{}) error

}

type pantiService struct {
	repo repository.Repository
}

func NewPantiService(repo repository.Repository) PantiService {
	return &pantiService{repo: repo}
}

func (s *pantiService) UpdatePanti(ctx context.Context, userID string, pantiID string, update map[string]interface{}) error {
	user, err := s.repo.User().GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user tidak ditemukan")
	}
	if user.PantiID == nil || user.PantiID.String() != pantiID {
		return fmt.Errorf("anda tidak memiliki akses ke panti ini")
	}
	panti, err := s.repo.Panti().GetByID(ctx, pantiID)
	if err != nil {
		return fmt.Errorf("panti tidak ditemukan")
	}
	if nama, ok := update["nama"].(string); ok {
		panti.Nama = nama
	}
	if alamat, ok := update["alamat"].(string); ok {
		panti.Alamat = &alamat
	}
	if telepon, ok := update["telepon"].(string); ok {
		panti.Telepon = &telepon
	}
	panti.UpdatedAt = time.Now()
	return s.repo.Panti().Update(ctx, panti)
}

func (s *pantiService) GetByID(ctx context.Context, id string) (*models.Panti, error) {
	return s.repo.Panti().GetByID(ctx, id)
}

func (s *pantiService) GetAll(ctx context.Context) ([]models.Panti, error) {
	return s.repo.Panti().GetAll(ctx)
}

func (s *pantiService) Create(ctx context.Context, panti *models.Panti) error {
	return s.repo.Panti().Create(ctx, panti)
}

func (s *pantiService) Update(ctx context.Context, panti *models.Panti) error {
	return s.repo.Panti().Update(ctx, panti)
}