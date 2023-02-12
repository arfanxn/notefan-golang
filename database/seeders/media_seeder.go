package seeders

import (
	"context"
	"database/sql"
	"notefan-golang/repositories"
	"runtime"
	"time"
)

type MediaSeeder struct {
	db        *sql.DB
	tableName string
	repo      *repositories.MediaRepo
}

func NewMediaSeeder(db *sql.DB) *MediaSeeder {
	return &MediaSeeder{
		db:        db,
		tableName: "medias",
		repo:      repositories.NewMediaRepo(db),
	}
}

func (seeder *MediaSeeder) Run() {
	// Consoler
	pc, _, _, _ := runtime.Caller(0)
	printStartRunningSeeder(pc)
	defer printFinishRunningSeeder(pc)

	// ---- Begin ----
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute/2) // Give a 30 second timeout
	defer cancel()

	_ = ctx

	// TODO : Complete media seeder
}
