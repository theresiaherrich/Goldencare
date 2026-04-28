package bootstrap

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/theresiaherrich/Goldencare/internal/config"
	"github.com/theresiaherrich/Goldencare/internal/database"
	"github.com/theresiaherrich/Goldencare/internal/repository"
	"github.com/theresiaherrich/Goldencare/internal/services"
)

type Container struct {
	Config     *config.Config
	DB         *sqlx.DB
	Repository repository.Repository
	Service    service.Service
}

func NewContainer() (*Container, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	log.Println("Configuration loaded")

	db, err := database.NewPostgres(cfg)
	if err != nil {
		return nil, err
	}
	log.Println("Database connected")

	repo := repository.NewPostgresRepository(db)
	log.Println("Repository initialized")
	svc := service.NewService(repo, cfg)

	log.Println("Services initialized")

	return &Container{
		Config:     cfg,
		DB:         db,
		Repository: repo,
		Service:    svc,
	}, nil
}

func (c *Container) Close() error {
	if c.Service != nil {
		if err := c.Service.Close(); err != nil {
			log.Printf("Error closing service: %v", err)
		}
	}
	if c.DB != nil {
		if err := c.DB.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}
	return nil
}
