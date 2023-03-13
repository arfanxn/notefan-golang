package entities

import (
	"github.com/google/uuid"
)

type NotificationUser struct {
	NotificationId uuid.UUID `json:"notification_id"`
	NotifierId     uuid.UUID `json:"notifier_id"`
	NotifiedId     uuid.UUID `json:"notified_id"`
}

/*
 * ----------------------------------------------------------------
 * NotificationUser Table and Columns methods  ⬇
 * ----------------------------------------------------------------
 */

// GetColumnNames returns the column names of the entity
func (ety NotificationUser) GetColumnNames() []string {
	return []string{
		"notification_id",
		"notifier_id",
		"notified_id",
	}
}

// GetTableName returns the table name
func (ety NotificationUser) GetTableName() string {
	return "notification_user"
}
