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

    existing, err := repo.Superadmin().GetByEmail(ctx, "admin@goldencare.com")
    if err == nil && existing != nil {
        log.Println("Superadmin already exists")
        return
    }

    hash, _ := bcrypt.GenerateFromPassword([]byte(cfg.SuperadminPassword), bcrypt.DefaultCost)

    superadmin := &models.Superadmin{
        ID:       uuid.New(),
        Name:     "Super Admin",
        Email:    "admin@goldencare.com",
        Password: string(hash),
    }

    if err := repo.Superadmin().Create(ctx, superadmin); err != nil {
        log.Fatalf("failed insert superadmin: %v", err)
    }
    log.Println("Superadmin created successfully")
}