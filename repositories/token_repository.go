package repositories

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/notefan-golang/helpers/stringh"
	"github.com/notefan-golang/models/entities"
	"github.com/notefan-golang/models/requests/query_reqs"
)

// TokenRepository represents a repository for token model/entity
type TokenRepository struct {
	db        *sql.DB
	Query     query_reqs.Query
	entity    entities.Token
	mutex     sync.Mutex
	waitGroup *sync.WaitGroup
}

// Instantiate a TokenRepository
func NewTokenRepository(db *sql.DB) *TokenRepository {
	return &TokenRepository{
		db:        db,
		Query:     query_reqs.Default(),
		entity:    entities.Token{},
		mutex:     sync.Mutex{},
		waitGroup: new(sync.WaitGroup),
	}
}

/*
 * ----------------------------------------------------------------
 * Repository utility methods ⬇
 * ----------------------------------------------------------------
 */

// scanRows scans rows of the database and returns it as structs, and returns error if any error has occurred.
func (repository *TokenRepository) scanRows(rows *sql.Rows) ([]entities.Token, error) {
	var tokens []entities.Token
	for rows.Next() {
		token := entities.Token{}
		err := rows.Scan(
			&token.Id,
			&token.TokenableType,
			&token.TokenableId,
			&token.Type,
			&token.Body,
			&token.UsedAt,
			&token.ExpiredAt,
			&token.CreatedAt,
			&token.UpdatedAt,
		)
		if err != nil {
			return tokens, err
		}
		tokens = append(tokens, token)
	}
	return tokens, nil
}

// scanRow scans only a row of the database and returns it as struct, and returns error if any error has occurred.
func (repository *TokenRepository) scanRow(rows *sql.Rows) (token entities.Token, err error) {
	tokens, err := repository.scanRows(rows)
	if err != nil {
		return
	}
	if len(tokens) == 0 {
		return
	}
	token, err = tokens[0], nil
	return
}

/*
 * ----------------------------------------------------------------
 * Repository query methods ⬇
 * ----------------------------------------------------------------
 */

// Find finds by id
func (repository *TokenRepository) Find(ctx context.Context, id string) (token entities.Token, err error) {
	query := "SELECT " + stringh.SliceColumnToStr(repository.entity.GetColumnNames()) +
		" FROM " + repository.entity.GetTableName() +
		" WHERE id = ?"
	rows, err := repository.db.QueryContext(ctx, query, id)
	if err != nil {
		return token, err
	}

	token, err = repository.scanRow(rows)
	if err != nil {
		return token, err
	}
	return token, nil
}

// FindByTokenableAndType finds a token by tokenable and type
func (repository *TokenRepository) FindByTokenableAndType(
	ctx context.Context, tokenableTyp string, tokenableId string, typ string,
) (token entities.Token, err error) {
	query := "SELECT " + stringh.SliceColumnToStr(repository.entity.GetColumnNames()) +
		" FROM " + repository.entity.GetTableName() +
		" WHERE `tokenable_type` = ? AND `tokenable_id` = ? AND `type` = ?"
	rows, err := repository.db.QueryContext(ctx, query, tokenableTyp, tokenableId, typ)
	if err != nil {
		return token, err
	}

	token, err = repository.scanRow(rows)
	if err != nil {
		return token, err
	}
	return token, nil
}

// Insert inserts tokens into the database
func (repository *TokenRepository) Insert(ctx context.Context, tokens ...*entities.Token) (sql.Result, error) {
	query := buildBatchInsertQuery(
		repository.entity.GetTableName(),
		len(tokens),
		repository.entity.GetColumnNames()...,
	)
	valueArgs := []any{}

	for _, token := range tokens {
		repository.waitGroup.Add(1)

		go func(wg *sync.WaitGroup, token *entities.Token) {

			defer wg.Done()

			repository.mutex.Lock()
			if token.Id == uuid.Nil {
				token.Id = uuid.New()
			}
			if token.CreatedAt.IsZero() {
				token.CreatedAt = time.Now()
			}
			valueArgs = append(valueArgs,
				token.Id,
				token.TokenableType,
				token.TokenableId,
				token.Type,
				token.Body,
				token.UsedAt,
				token.ExpiredAt,
				token.CreatedAt,
				token.UpdatedAt,
			)
			repository.mutex.Unlock()
		}(repository.waitGroup, token)
	}

	repository.waitGroup.Wait()

	result, err := repository.db.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		return result, err
	}
	return result, nil
}

// Create creates token into database
func (repository *TokenRepository) Create(ctx context.Context, token *entities.Token) (sql.Result, error) {
	result, err := repository.Insert(ctx, token)
	if err != nil {
		return result, err
	}

	return result, nil
}

// UpdateById
func (repository *TokenRepository) UpdateById(ctx context.Context, token *entities.Token) (sql.Result, error) {
	query := buildUpdateQuery(repository.entity.GetTableName(), repository.entity.GetColumnNames()...) +
		" WHERE id = ?"

	// Refresh entity updated at
	token.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}

	result, err := repository.db.ExecContext(ctx, query,
		token.Id,
		token.TokenableType,
		token.TokenableId,
		token.Type,
		token.Body,
		token.UsedAt,
		token.ExpiredAt,
		token.CreatedAt,
		token.UpdatedAt,
		token.Id,
	)

	return result, err
}

// DeleteByIds deletes the entities associated with the given ids
func (repository *TokenRepository) DeleteByIds(ctx context.Context, ids ...string) (sql.Result, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	query, valueArgs := buildBatchDeleteQueryByIds(repository.entity.GetTableName(), ids...)
	result, err := repository.db.ExecContext(ctx, query, valueArgs...)
	return result, err
}
