package seeders

import (
	"database/sql"
	"notefan-golang/database/factories"
	"notefan-golang/helper"
	"notefan-golang/repositories"
	"runtime"
)

type SpaceSeeder struct {
	db        *sql.DB
	tableName string
	repo      *repositories.SpaceRepo
}

func NewSpaceSeeder(db *sql.DB) *SpaceSeeder {
	return &SpaceSeeder{
		db:        db,
		tableName: "spaces",
		repo:      repositories.NewSpaceRepo(db),
	}
}

func (seeder *SpaceSeeder) Run() {
	// Consoler
	pc, _, _, _ := runtime.Caller(0)
	printStartRunning(pc)
	defer printFinishRunning(pc)

	// ---- Begin ----
	totalRows := 50
	valueArgs := []any{}

	for i := 0; i < totalRows; i++ {
		space := factories.NewSpace()
		valueArgs = append(
			valueArgs,
			space.Id.String(), space.Name, space.Description, space.Domain, space.CreatedAt, space.UpdatedAt)
	}

	query := helper.BuildBulkInsertQuery(seeder.tableName, totalRows,
		`id`, `name`, `description`, `domain`, `created_at`, `updated_at`)

	stmt, err := seeder.db.Prepare(query)
	helper.LogFatalIfError(err)

	_, err = stmt.Exec(valueArgs...)
	helper.LogFatalIfError(err)
}
