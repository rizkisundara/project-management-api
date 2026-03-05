package models

import (
	"github.com/google/uuid"
	"github.com/rizkisundara/project-management-api/models/types"
)

type CardPosition struct {
	InternalID int64           `json:"internal_id" db:"internal_id" gorm:"primaryKey;autoIncrement"`
	PublicID   uuid.UUID       `json:"public_id" db:"public_id"`
	ListID     int64           `json:"list_internal_id" db:"list_internal_id" gorm:"column:list_internal_id"`
	CardOrder  types.UUIDArray `json:"card_order" db:"card_order" gorm:"type:uuid[]"`
}
