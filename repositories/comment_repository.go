package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/notefan-golang/helpers/reflecth"
	"github.com/notefan-golang/helpers/stringh"
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
		columnNames: reflecth.GetFieldJsonTag(entities.Comment{}),
	}
}

func (repository *CommentRepository) All(ctx context.Context) (comments []entities.Comment, err error) {
	query := "SELECT " + stringh.SliceColumnToStr(repository.columnNames) + " FROM " + repository.tableName
	rows, err := repository.db.QueryContext(ctx, query)
	if err != nil {
		return
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
			return comments, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (repository *CommentRepository) Insert(ctx context.Context, comments ...*entities.Comment) (
	sql.Result, error) {
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
	result, err := repository.db.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (repository *CommentRepository) Create(ctx context.Context, comment *entities.Comment) (sql.Result, error) {
	return repository.Insert(ctx, comment)
}
