package database

import (
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib" 
	"github.com/jmoiron/sqlx"
	"github.com/theresiaherrich/Goldencare/internal/config"
)

func NewPostgres(cfg *config.Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", cfg.DBURL)
	if err != nil {
		return nil, fmt.Errorf("gagal membuka koneksi database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("database tidak bisa dijangkau: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)

	log.Println("Database PostgreSQL terhubung!")
	return db, nil
}
