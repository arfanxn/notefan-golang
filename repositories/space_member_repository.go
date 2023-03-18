package repositories

import (
	"bytes"
	"context"
	"database/sql"
	"strings"

	"github.com/notefan-golang/helpers/entityh"
	"github.com/notefan-golang/helpers/stringh"
	"github.com/notefan-golang/models/entities"
	"github.com/notefan-golang/models/requests/query_reqs"
)

type SpaceMemberRepository struct {
	db     *sql.DB
	Query  query_reqs.Query
	entity entities.User
}

func NewSpaceMemberRepository(db *sql.DB) *SpaceMemberRepository {
	return &SpaceMemberRepository{
		db:     db,
		Query:  query_reqs.Default(),
		entity: entities.User{},
	}
}

/*
 * ----------------------------------------------------------------
 * Repository utilty methods ⬇
 * ----------------------------------------------------------------
 */

//

/*
 * ----------------------------------------------------------------
 * Repository query methods ⬇
 * ----------------------------------------------------------------
 */

// GetBySpaceId returns the members of the given space id
func (repository *SpaceMemberRepository) GetBySpaceId(ctx context.Context, spaceId string) (
	members []entities.User, err error) {

	var (
		valueArgs []any
	)

	queryBuf := bytes.NewBufferString("SELECT ")
	queryBuf.WriteString(stringh.SliceTableColumnToStr(
		repository.entity.GetTableName(), repository.entity.GetColumnNames()),
	) // column names
	queryBuf.WriteRune(',') // comma user_role_space table columns
	queryBuf.WriteString(stringh.SliceTableColumnToStr(
		entityh.GetTableName(entities.UserRoleSpace{}),
		entityh.GetColumnNames(entities.UserRoleSpace{}),
	))
	queryBuf.WriteRune(',') // comma role table columns
	queryBuf.WriteString(stringh.SliceTableColumnToStr(
		entityh.GetTableName(entities.Role{}),
		entityh.GetColumnNames(entities.Role{}),
	))
	queryBuf.WriteString(" FROM ")
	queryBuf.WriteString(repository.entity.GetTableName())
	queryBuf.WriteString(" INNER JOIN ") // inner join UserRoleSpace relation
	queryBuf.WriteString(entityh.GetTableName(entities.UserRoleSpace{}))
	queryBuf.WriteString(" ON ")
	queryBuf.WriteString(repository.entity.GetTableName() + ".`id`") // "user.`id`"
	queryBuf.WriteString(" = ")
	queryBuf.WriteString(entityh.GetTableName(entities.UserRoleSpace{}) + ".user_id")
	queryBuf.WriteString(" INNER JOIN ") // inner join Role relation
	queryBuf.WriteString(entityh.GetTableName(entities.Role{}))
	queryBuf.WriteString(" ON ")
	queryBuf.WriteString(entityh.GetTableName(entities.UserRoleSpace{}) + ".`role_id`") // "user.`id`"
	queryBuf.WriteString(" = ")
	queryBuf.WriteString(entityh.GetTableName(entities.Role{}) + ".id")
	queryBuf.WriteString(" WHERE ")
	queryBuf.WriteString(entityh.GetTableName(entities.UserRoleSpace{}) + ".`space_id` = ?")
	valueArgs = append(valueArgs, spaceId)
	if len(repository.Query.OrderBys) != 0 {
		queryBuf.WriteString(" ORDER BY ")
		index := 0
		for columnName, orderingType := range repository.Query.OrderBys {
			if index > 0 {
				queryBuf.WriteRune(',')
			}
			queryBuf.WriteString(" " + repository.entity.GetTableName() + "." + columnName + " " + strings.ToUpper(orderingType) + " ")
			index++
		}
	}
	queryBuf.WriteString(" LIMIT ? OFFSET ? ")
	valueArgs = append(valueArgs, repository.Query.Limit, repository.Query.Offset)
	if err != nil {
		return
	}

	rows, err := repository.db.QueryContext(ctx, queryBuf.String(), valueArgs...)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var user entities.User
		var urs entities.UserRoleSpace
		err := rows.Scan(
			&user.Id,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.CreatedAt,
			&user.UpdatedAt,
			&urs.Id,
			&urs.UserId,
			&urs.RoleId,
			&urs.SpaceId,
			&urs.CreatedAt,
			&urs.UpdatedAt,
			&user.Role.Id,
			&user.Role.Name,
		)
		if err != nil {
			return members, err
		}
		members = append(members, user)
	}
	return members, nil
}

// FindByMemberIdAndSpaceId finds member/user by member id and space id
func (repository *SpaceMemberRepository) FindByMemberIdAndSpaceId(ctx context.Context, memberId string, spaceId string) (
	member entities.User, err error) {

	var (
		valueArgs []any
	)

	queryBuf := bytes.NewBufferString("SELECT ")
	queryBuf.WriteString(stringh.SliceTableColumnToStr(
		repository.entity.GetTableName(), repository.entity.GetColumnNames()),
	) // column names
	queryBuf.WriteRune(',') // comma user_role_space table columns
	queryBuf.WriteString(stringh.SliceTableColumnToStr(
		entityh.GetTableName(entities.UserRoleSpace{}),
		entityh.GetColumnNames(entities.UserRoleSpace{}),
	))
	queryBuf.WriteRune(',') // comma role table columns
	queryBuf.WriteString(stringh.SliceTableColumnToStr(
		entityh.GetTableName(entities.Role{}),
		entityh.GetColumnNames(entities.Role{}),
	))
	queryBuf.WriteString(" FROM ")
	queryBuf.WriteString(repository.entity.GetTableName())
	queryBuf.WriteString(" INNER JOIN ") // inner join UserRoleSpace relation
	queryBuf.WriteString(entityh.GetTableName(entities.UserRoleSpace{}))
	queryBuf.WriteString(" ON ")
	queryBuf.WriteString(repository.entity.GetTableName() + ".`id`") // "user.`id`"
	queryBuf.WriteString(" = ")
	queryBuf.WriteString(entityh.GetTableName(entities.UserRoleSpace{}) + ".user_id")
	queryBuf.WriteString(" INNER JOIN ") // inner join Role relation
	queryBuf.WriteString(entityh.GetTableName(entities.Role{}))
	queryBuf.WriteString(" ON ")
	queryBuf.WriteString(entityh.GetTableName(entities.UserRoleSpace{}) + ".`role_id`") // "user.`id`"
	queryBuf.WriteString(" = ")
	queryBuf.WriteString(entityh.GetTableName(entities.Role{}) + ".id")
	queryBuf.WriteString(" WHERE ")
	queryBuf.WriteString(repository.entity.GetTableName() + ".`id` = ?") // where users.`id` = ?
	valueArgs = append(valueArgs, memberId)
	queryBuf.WriteString(" AND ")
	queryBuf.WriteString(entityh.GetTableName(entities.UserRoleSpace{}) + ".`space_id` = ?") // where user_role_space.`space_id` = ?
	valueArgs = append(valueArgs, spaceId)
	queryBuf.WriteString(" LIMIT ? OFFSET ? ")
	valueArgs = append(valueArgs, repository.Query.Limit, repository.Query.Offset)
	if err != nil {
		return
	}

	rows, err := repository.db.QueryContext(ctx, queryBuf.String(), valueArgs...)
	if err != nil {
		return
	}
	defer rows.Close()
	if rows.Next() {
		var urs entities.UserRoleSpace
		err = rows.Scan(
			&member.Id,
			&member.Name,
			&member.Email,
			&member.Password,
			&member.CreatedAt,
			&member.UpdatedAt,
			&urs.Id,
			&urs.UserId,
			&urs.RoleId,
			&urs.SpaceId,
			&urs.CreatedAt,
			&urs.UpdatedAt,
			&member.Role.Id,
			&member.Role.Name,
		)
		if err != nil {
			return
		}
		return member, nil
	} else {
		return member, err
	}
}
