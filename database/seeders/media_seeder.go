package seeders

import (
	"context"
	"database/sql"
	"notefan-golang/repositories"
	"time"
)

type MediaSeeder struct {
	db   *sql.DB
	repo *repositories.MediaRepo
}

func NewMediaSeeder(db *sql.DB) *MediaSeeder {
	return &MediaSeeder{
		db:   db,
		repo: repositories.NewMediaRepo(db),
	}
}

func (seeder *MediaSeeder) Run() {

	// ---- Begin ----
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute/2) // Give a 30 second timeout
	defer cancel()

	_ = ctx

	// TODO : Complete media seeder
}
