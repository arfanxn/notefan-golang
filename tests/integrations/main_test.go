package integrations

import (
	"fmt"
	"os"
	"testing"

	"github.com/notefan-golang/helper"

	"github.com/golang-migrate/migrate/v4"
)

func TestMain(m *testing.M) {
	setup()
	m.Run()
	teardown()
}

// setup sets up the test
func setup() {
	migrateUpTestDBTables()
}

// teardown teardowns the test
func teardown() {
	migrateDownTestDBTables()
}

func migrateUpTestDBTables() {
	fmt.Println("Working directory ")
	fmt.Println(os.Getwd())
	m, err := migrate.New("github.com/notefan-golang/database/migrations", "mysql://root@tcp(localhost:3306)/notefan_test")
	helper.ErrorPanic(err)
	m.Up()
}

func migrateDownTestDBTables() {
	m, err := migrate.New("github.com/notefan-golang/database/migrations", "mysql://root@tcp(localhost:3306)/notefan_test")
	helper.ErrorPanic(err)
	m.Drop()
}
