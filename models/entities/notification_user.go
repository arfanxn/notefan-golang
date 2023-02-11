package entities

import (
	"github.com/google/uuid"
)

type NotificationUser struct {
	NotificationId uuid.UUID `json:"notification_id"`
	NotifierId     uuid.UUID `json:"notifier_id"`
	NotifiedId     uuid.UUID `json:"notified_id"`
}
