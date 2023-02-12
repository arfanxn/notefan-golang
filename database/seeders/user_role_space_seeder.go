package seeders

import (
	"context"
	"database/sql"
	"math/rand"
	"notefan-golang/helper"
	"notefan-golang/models/entities"
	"notefan-golang/repositories"
	"runtime"
	"time"
)

type UserRoleSpaceSeeder struct {
	db        *sql.DB
	tableName string
	repo      repositories.UserRoleSpaceRepo
}

func NewUserRoleSpaceSeeder(db *sql.DB) *UserRoleSpaceSeeder {
	return &UserRoleSpaceSeeder{
		db:        db,
		tableName: "user_role_space",
		repo:      *repositories.NewUserRoleSpaceRepo(db),
	}
}

func (seeder *UserRoleSpaceSeeder) Run() {
	// Consoler
	pc, _, _, _ := runtime.Caller(0)
	printStartRunning(pc)
	defer printFinishRunning(pc)

	// ---- Begin ----
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	userRepo := repositories.NewUserRepo(seeder.db)
	users, err := userRepo.All(ctx)
	helper.PanicIfError(err)

	roleRepo := repositories.NewRoleRepo(seeder.db)
	roleSpaceOwner, err := roleRepo.FindByName(ctx, "space owner")
	helper.PanicIfError(err)
	roleSpaceMember, err := roleRepo.FindByName(ctx, "space member")
	helper.PanicIfError(err)

	spaceRepo := repositories.NewSpaceRepo(seeder.db)
	spaces, err := spaceRepo.All(ctx)
	helper.PanicIfError(err)

	totalRows := len(users)
	valueArgs := []any{}

	for _, user := range users {
		role := entities.Role{}
		if helper.BooleanRandom() {
			role = roleSpaceOwner
		} else {
			role = roleSpaceMember
		}

		space := spaces[rand.Intn(len(spaces))]

		urs := entities.UserRoleSpace{
			UserId:    user.Id,
			RoleId:    role.Id,
			SpaceId:   space.Id,
			CreatedAt: time.Now(),
			UpdatedAt: sql.NullTime{
				Time:  time.Now(),
				Valid: true},
		}

		if user.Email == "arfan@gmail.com" {
			urs = entities.UserRoleSpace{
				UserId:    user.Id,
				RoleId:    roleSpaceOwner.Id,
				SpaceId:   space.Id,
				CreatedAt: time.Now(),
				UpdatedAt: sql.NullTime{Time: time.Now().AddDate(0, 0, 1), Valid: true},
			}
		}

		valueArgs = append(
			valueArgs,
			urs.UserId.String(), urs.RoleId.String(), urs.SpaceId.String(), urs.CreatedAt, urs.UpdatedAt)
	}

	query := helper.BuildBulkInsertQuery(seeder.tableName, totalRows,
		`user_id`, `role_id`, `space_id`, `created_at`, `updated_at`)

	stmt, err := seeder.db.Prepare(query)
	helper.PanicIfError(err)

	_, err = stmt.Exec(valueArgs...)
	helper.PanicIfError(err)
}
