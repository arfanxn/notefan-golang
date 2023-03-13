package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/notefan-golang/helpers/stringh"
	"github.com/notefan-golang/models/entities"
	"github.com/notefan-golang/models/requests/query_reqs"

	"github.com/google/uuid"
)

type CommentReactionRepository struct {
	db     *sql.DB
	Query  query_reqs.Query
	entity entities.CommentReaction
}

func NewCommentReactionRepository(db *sql.DB) *CommentReactionRepository {
	return &CommentReactionRepository{
		db:     db,
		Query:  query_reqs.Default(),
		entity: entities.CommentReaction{},
	}
}

func (repository *CommentReactionRepository) All(ctx context.Context) (
	commentReactions []entities.CommentReaction, err error) {
	query := "SELECT " + stringh.SliceColumnToStr(repository.entity.GetColumnNames()) +
		" FROM " + repository.entity.GetTableName()
	rows, err := repository.db.QueryContext(ctx, query)
	if err != nil {
		return
	}

	for rows.Next() {
		commentReaction := entities.CommentReaction{}
		err := rows.Scan(
			&commentReaction.Id,
			&commentReaction.CommentId,
			&commentReaction.UserId,
			&commentReaction.CreatedAt,
			&commentReaction.UpdatedAt,
		)
		if err != nil {
			return commentReactions, err
		}
		commentReactions = append(commentReactions, commentReaction)
	}
	return commentReactions, nil
}

func (repository *CommentReactionRepository) Insert(ctx context.Context, commentReactions ...*entities.CommentReaction) (sql.Result, error) {
	query := buildBatchInsertQuery(
		repository.entity.GetTableName(),
		len(commentReactions),
		repository.entity.GetColumnNames()...,
	)
	valueArgs := []any{}

	for _, commentReaction := range commentReactions {
		if commentReaction.Id == uuid.Nil {
			commentReaction.Id = uuid.New()
		}
		if commentReaction.CreatedAt.IsZero() {
			commentReaction.CreatedAt = time.Now()
		}
		valueArgs = append(valueArgs,
			commentReaction.Id,
			commentReaction.CommentId,
			commentReaction.UserId,
			commentReaction.CreatedAt,
			commentReaction.UpdatedAt,
		)
	}

	result, err := repository.db.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (repository *CommentReactionRepository) Create(ctx context.Context, commentReaction *entities.CommentReaction) (sql.Result, error) {
	return repository.Insert(ctx, commentReaction)
}
