package entities

type EntityContract interface {
	GetTableName() string
	GetColumnNames() []string
}
