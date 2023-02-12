package entities

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type PageContentChangeHistory struct {
	BeforePageContentChangeHistory uuid.UUID    `json:"before_page_content_id"`
	AfterPageContentChangeHistory  uuid.UUID    `json:"after_page_content_id"`
	UserId                         uuid.UUID    `json:"user_id"`
	CreatedAt                      time.Time    `json:"created_at"`
	UpdatedAt                      sql.NullTime `json:"updated_at"`
}