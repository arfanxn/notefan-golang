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

type CommentReactionRepo struct {
	db          *sql.DB
	tableName   string
	columnNames []string
}

func NewCommentReactionRepo(db *sql.DB) *CommentReactionRepo {
	return &CommentReactionRepo{
		db:          db,
		tableName:   "comment_reactions",
		columnNames: helper.GetStructFieldJsonTag(entities.CommentReaction{}),
	}
}

func (repo *CommentReactionRepo) All(ctx context.Context) ([]entities.CommentReaction, error) {
	query := "SELECT " + helper.SliceTableColumnsToString(repo.columnNames) + " FROM " + repo.tableName
	commentReactions := []entities.CommentReaction{}
	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		helper.LogIfError(err)
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
			helper.LogIfError(err)
			return commentReactions, err
		}
		commentReactions = append(commentReactions, commentReaction)
	}

	if len(commentReactions) == 0 {
		return commentReactions, exceptions.DataNotFoundError
	}

	return commentReactions, nil
}

func (repo *CommentReactionRepo) Insert(ctx context.Context, commentReactions ...entities.CommentReaction) ([]entities.CommentReaction, error) {
	query := buildBatchInsertQuery(repo.tableName, len(commentReactions), repo.columnNames...)
	valueArgs := []any{}

	for _, commentReaction := range commentReactions {
		if commentReaction.Id.String() == "" {
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

	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		helper.LogIfError(err)
		return commentReactions, err
	}
	_, err = stmt.ExecContext(ctx, valueArgs...)
	if err != nil {
		helper.LogIfError(err)
		return commentReactions, err
	}
	return commentReactions, nil
}

func (repo *CommentReactionRepo) Create(ctx context.Context, commentReaction entities.CommentReaction) (entities.CommentReaction, error) {
	commentReactions, err := repo.Insert(ctx, commentReaction)
	if err != nil {
		return entities.CommentReaction{}, err
	}

	return commentReactions[0], nil
}
