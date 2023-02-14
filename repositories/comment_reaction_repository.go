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

type CommentReactionRepository struct {
	db          *sql.DB
	tableName   string
	columnNames []string
}

func NewCommentReactionRepository(db *sql.DB) *CommentReactionRepository {
	return &CommentReactionRepository{
		db:          db,
		tableName:   "comment_reactions",
		columnNames: helper.ReflectGetStructFieldJsonTag(entities.CommentReaction{}),
	}
}

func (repository *CommentReactionRepository) All(ctx context.Context) ([]entities.CommentReaction, error) {
	query := "SELECT " + helper.DBSliceColumnsToStr(repository.columnNames) + " FROM " + repository.tableName
	commentReactions := []entities.CommentReaction{}
	rows, err := repository.db.QueryContext(ctx, query)
	if err != nil {
		helper.ErrorLog(err)
		return commentReactions, err
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
			helper.ErrorLog(err)
			return commentReactions, err
		}
		commentReactions = append(commentReactions, commentReaction)
	}

	if len(commentReactions) == 0 {
		return commentReactions, exceptions.HTTPNotFound
	}

	return commentReactions, nil
}

func (repository *CommentReactionRepository) Insert(ctx context.Context, commentReactions ...entities.CommentReaction) ([]entities.CommentReaction, error) {
	query := buildBatchInsertQuery(repository.tableName, len(commentReactions), repository.columnNames...)
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

	stmt, err := repository.db.PrepareContext(ctx, query)
	if err != nil {
		helper.ErrorLog(err)
		return commentReactions, err
	}
	_, err = stmt.ExecContext(ctx, valueArgs...)
	if err != nil {
		helper.ErrorLog(err)
		return commentReactions, err
	}
	return commentReactions, nil
}

func (repository *CommentReactionRepository) Create(ctx context.Context, commentReaction entities.CommentReaction) (entities.CommentReaction, error) {
	commentReactions, err := repository.Insert(ctx, commentReaction)
	if err != nil {
		return entities.CommentReaction{}, err
	}

	return commentReactions[0], nil
}
