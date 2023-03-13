package entityh

import "github.com/notefan-golang/models/entities"

// GetTableName returns table name of the entity
func GetTableName[T entities.EntityContract](entity T) string {
	return entity.GetTableName()
}

// GetColumnNames returns column names of the entity
func GetColumnNames[T entities.EntityContract](entity T) []string {
	return entity.GetColumnNames()
}
