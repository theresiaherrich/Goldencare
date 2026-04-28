package service

import (
	"context"

	"github.com/theresiaherrich/Goldencare/internal/models"
	"github.com/theresiaherrich/Goldencare/internal/repository"
)

type UserService interface {
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id string) error
}

type userService struct {
	repo repository.Repository
}

func NewUserService(repo repository.Repository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetByID(ctx context.Context, id string) (*models.User, error) {
	return s.repo.User().GetByID(ctx, id)
}

func (s *userService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.repo.User().GetByEmail(ctx, email)
}

func (s *userService) Update(ctx context.Context, user *models.User) error {
	return s.repo.User().Update(ctx, user)
}

func (s *userService) Delete(ctx context.Context, id string) error {
	return s.repo.User().Delete(ctx, id)
}