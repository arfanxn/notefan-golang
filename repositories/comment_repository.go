package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/notefan-golang/exceptions"
	"github.com/notefan-golang/helper"
	"github.com/notefan-golang/models/entities"

	"github.com/google/uuid"
)

type CommentRepository struct {
	db          *sql.DB
	tableName   string
	columnNames []string
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{
		db:          db,
		tableName:   "comments",
		columnNames: helper.ReflectGetStructFieldJsonTag(entities.Comment{}),
	}
}

func (repository *CommentRepository) All(ctx context.Context) ([]entities.Comment, error) {
	query := "SELECT " + helper.DBSliceColumnsToStr(repository.columnNames) + " FROM " + repository.tableName
	comments := []entities.Comment{}
	rows, err := repository.db.QueryContext(ctx, query)
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

func (repository *CommentRepository) Insert(ctx context.Context, comments ...entities.Comment) ([]entities.Comment, error) {
	query := buildBatchInsertQuery(repository.tableName, len(comments), repository.columnNames...)
	valueArgs := []any{}

	for _, comment := range comments {
		if comment.Id == uuid.Nil {
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

	stmt, err := repository.db.PrepareContext(ctx, query)
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

func (repository *CommentRepository) Create(ctx context.Context, comment entities.Comment) (entities.Comment, error) {
	comments, err := repository.Insert(ctx, comment)
	if err != nil {
		return entities.Comment{}, err
	}

	return comments[0], nil
}
