package seeders

import (
	"context"
	"database/sql"
	"math/rand"
	"notefan-golang/database/factories"
	"notefan-golang/helper"
	"notefan-golang/repositories"
	"runtime"
	"time"
)

type PageSeeder struct {
	db        *sql.DB
	tableName string
	repo      *repositories.PageRepo
	spaceRepo *repositories.SpaceRepo
}

func NewPageSeeder(db *sql.DB) *PageSeeder {
	return &PageSeeder{
		db:        db,
		tableName: "pages",
		repo:      repositories.NewPageRepo(db),
		spaceRepo: repositories.NewSpaceRepo(db),
	}
}

func (seeder *PageSeeder) Run() {
	// Consoler
	pc, _, _, _ := runtime.Caller(0)
	printStartRunning(pc)
	defer printFinishRunning(pc)

	// ---- Begin ----
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	spaces, err := seeder.spaceRepo.All(ctx)

	totalRows := len(spaces) * 2
	valueArgs := []any{}

	for i := 0; i < totalRows; i++ {
		space := spaces[rand.Intn(len(spaces))]
		page := factories.NewPage()
		page.SpaceId = space.Id
		valueArgs = append(
			valueArgs,
			page.Id.String(), page.SpaceId.String(), page.Title, page.Order, page.CreatedAt, page.UpdatedAt)
	}

	query := helper.BuildBulkInsertQuery(seeder.tableName, totalRows,
		`id`, `space_id`, `title`, `order`, `created_at`, `updated_at`)

	stmt, err := seeder.db.Prepare(query)
	helper.PanicIfError(err)

	_, err = stmt.Exec(valueArgs...)
	helper.PanicIfError(err)
}
