package repositories

import "database/sql"

type UserRoleSpaceRepo struct {
	tableName string
	db        *sql.DB
}

func NewUserRoleSpaceRepo(db *sql.DB) *UserRoleSpaceRepo {
	return &UserRoleSpaceRepo{
		tableName: "user_role_space",
		db:        db,
	}
}

func (repo *UserRoleSpaceRepo) FindByName() {
	//
}
