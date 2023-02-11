package seeders

import (
	"notefan-golang/database/factories"
	"notefan-golang/helper"
	"runtime"
)

func UserSeeder(s *Seeder) {
	// Consoler
	pc, _, _, _ := runtime.Caller(0)
	printStartRunning(pc)
	defer printFinishRunning(pc)

	// Begin
	totalRows := 50
	valueArgs := []any{}
	for i := 0; i < totalRows; i++ {
		user := factories.NewUser()
		valueArgs = append(
			valueArgs,
			user.Id.String(), user.Name, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
	}

	query := helper.BuildBulkInsertQuery("users", totalRows,
		`id`, `name`, `email`, `password`, `created_at`, `updated_at`)

	stmt, err := s.db.Prepare(query)
	helper.LogFatalIfError(err)

	_, err = stmt.Exec(valueArgs...)
	helper.LogFatalIfError(err)
}
