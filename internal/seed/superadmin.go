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

	// cek existing superadmin
	user, err := repo.User().FindByEmail(ctx, "admin@goldencare.com")
	if err == nil && user != nil {
		log.Println("Superadmin already exists")
		return
	}

	// WAJIB dari ENV
	pass := cfg.SuperadminPassword
	if pass == "" {
		log.Fatal("SUPERADMIN_PASSWORD is not set in environment")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("bcrypt error: %v", err)
	}

	superadmin := &models.User{
		ID:       uuid.New(),
		Name:     "Super Admin",
		Email:    "admin@goldencare.com",
		Password: string(hash),
		Role:     "superadmin",
	}

	if err := repo.User().Create(ctx, superadmin); err != nil {
		log.Fatalf("failed insert superadmin: %v", err)
	}

	log.Println("Superadmin CREATED FROM ENV SUCCESSFULLY")
}