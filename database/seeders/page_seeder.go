package seeders

import (
	"runtime"
)

func PageSeeder() {
	// Consoler
	pc, _, _, _ := runtime.Caller(0)
	printStartRunning(pc)
	defer printFinishRunning(pc)

	// ---- Begin ----
	// tableName := "pages"
	// totalRows := 50
	// valueArgs := []any{}

	// for i := 0; i < totalRows; i++ {
	// 	space := factories.NewSpace()
	// 	valueArgs = append(
	// 		valueArgs,
	// 		space.Id.String(), space.Name, space.Description, space.Domain, space.CreatedAt, space.UpdatedAt)
	// }

	// query := helper.BuildBulkInsertQuery(tableName, totalRows,
	// 	`id`, `name`, `description`, `domain`, `created_at`, `updated_at`)

	// stmt, err := seeder.db.Prepare(query)
	// helper.LogFatalIfError(err)

	// _, err = stmt.Exec(valueArgs...)
	// helper.LogFatalIfError(err)
}
