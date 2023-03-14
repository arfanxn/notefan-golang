package entities

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type PageContentChangeHistory struct {
	BeforePageContentId uuid.UUID    `json:"before_page_content_id"`
	AfterPageContentId  uuid.UUID    `json:"after_page_content_id"`
	UserId              uuid.UUID    `json:"user_id"`
	CreatedAt           time.Time    `json:"created_at"`
	UpdatedAt           sql.NullTime `json:"updated_at"`

	// Relations

	BeforePageContent PageContent `json:"before_page_content"`
	AfterPageContent  PageContent `json:"after_page_content"`
	User              User        `json:"user"`
}

/*
 * ----------------------------------------------------------------
 * PageContentChangeHistory Table and Columns methods  ⬇
 * ----------------------------------------------------------------
 */

// GetColumnNames returns the column names of the entity
func (ety PageContentChangeHistory) GetColumnNames() []string {
	return []string{
		"before_page_content_id",
		"after_page_content_id",
		"user_id",
		"created_at",
		"updated_at",
	}
}

// GetTableName returns the table name
func (ety PageContentChangeHistory) GetTableName() string {
	return "page_content_change_history"
}
