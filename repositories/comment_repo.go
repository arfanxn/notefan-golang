package repositories

import (
	"context"
	"database/sql"
	"notefan-golang/exceptions"
	"notefan-golang/helper"
	"notefan-golang/models/entities"
	"time"

	"github.com/google/uuid"
)

type CommentRepo struct {
	db          *sql.DB
	tableName   string
	columnNames []string
}

func NewCommentRepo(db *sql.DB) *CommentRepo {
	return &CommentRepo{
		db:          db,
		tableName:   "comments",
		columnNames: helper.ReflectGetStructFieldJsonTag(entities.Comment{}),
	}
}

func (repo *CommentRepo) All(ctx context.Context) ([]entities.Comment, error) {
	query := "SELECT " + helper.DBSliceColumnsToStr(repo.columnNames) + " FROM " + repo.tableName
	comments := []entities.Comment{}
	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		helper.ErrorLog(err)
		return comments, err
	}

	for rows.Next() {
		comment := entities.Comment{}
		err := rows.Scan(
			&comment.Id,
			&comment.CommentableType,
			&comment.CommentableId,
			&comment.UserId,
			&comment.Body,
			&comment.ResolvedAt,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		)
		if err != nil {
			helper.ErrorLog(err)
			return comments, err
		}
		comments = append(comments, comment)
	}

	if len(comments) == 0 {
		return comments, exceptions.HTTPNotFound
	}

	return comments, nil
}

func (repo *CommentRepo) Insert(ctx context.Context, comments ...entities.Comment) ([]entities.Comment, error) {
	query := buildBatchInsertQuery(repo.tableName, len(comments), repo.columnNames...)
	valueArgs := []any{}

	for _, comment := range comments {
		if comment.Id.String() == "" {
			comment.Id = uuid.New()
		}
		if comment.CreatedAt.IsZero() {
			comment.CreatedAt = time.Now()
		}
		valueArgs = append(valueArgs,
			comment.Id,
			comment.CommentableType,
			comment.CommentableId,
			comment.UserId,
			comment.Body,
			comment.ResolvedAt,
			comment.CreatedAt,
			comment.UpdatedAt,
		)
	}

	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		helper.ErrorLog(err)
		return comments, err
	}
	_, err = stmt.ExecContext(ctx, valueArgs...)
	if err != nil {
		helper.ErrorLog(err)
		return comments, err
	}
	return comments, nil
}

func (repo *CommentRepo) Create(ctx context.Context, comment entities.Comment) (entities.Comment, error) {
	comments, err := repo.Insert(ctx, comment)
	if err != nil {
		return entities.Comment{}, err
	}

	return comments[0], nil
}
