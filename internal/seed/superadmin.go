package seed

import (
	"context"
	"log"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
	"github.com/theresiaherrich/Goldencare/internal/config"
	"github.com/theresiaherrich/Goldencare/internal/models"
	"github.com/theresiaherrich/Goldencare/internal/repository"
)

func SeedSuperadmin(repo repository.Repository, cfg *config.Config) {
	ctx := context.Background()

	existing, err := repo.User().FindSuperadmin(ctx)
	if err != nil {
		log.Println("error check superadmin:", err)
		return
	}

	if existing != nil {
		log.Println("superadmin already exists, skipping seed")
		return
	}

	if cfg.SuperadminPassword == "" {
		log.Fatal("SUPERADMIN_PASSWORD is not set in env")
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(cfg.SuperadminPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		log.Fatal("failed to hash password:", err)
	}

	superadmin := &models.User{
		ID:       uuid.New(),
		Name:     "Super Admin",
		Email:    "admin@goldencare.com",
		Password: string(passwordHash),
		Role:     "superadmin",
	}

	if err := repo.User().Create(ctx, superadmin); err != nil {
		log.Fatal("failed to create superadmin:", err)
	}

	log.Println("superadmin seeded successfully")
}