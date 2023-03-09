package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/notefan-golang/helpers/errorh"
	"github.com/notefan-golang/helpers/reflecth"
	"github.com/notefan-golang/helpers/stringh"
	"github.com/notefan-golang/models/entities"

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
		columnNames: reflecth.GetFieldJsonTag(entities.CommentReaction{}),
	}
}

func (repository *CommentReactionRepository) All(ctx context.Context) (
	commentReactions []entities.CommentReaction, err error) {
	query := "SELECT " + stringh.SliceColumnToStr(repository.columnNames) + " FROM " + repository.tableName
	rows, err := repository.db.QueryContext(ctx, query)
	if err != nil {
		errorh.Log(err)
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
			errorh.Log(err)
			return commentReactions, err
		}
		commentReactions = append(commentReactions, commentReaction)
	}
	return commentReactions, nil
}

func (repository *CommentReactionRepository) Insert(ctx context.Context, commentReactions ...*entities.CommentReaction) (sql.Result, error) {
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

	result, err := repository.db.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		errorh.Log(err)
		return result, err
	}
	return result, nil
}

func (repository *CommentReactionRepository) Create(ctx context.Context, commentReaction *entities.CommentReaction) (sql.Result, error) {
	return repository.Insert(ctx, commentReaction)
}
