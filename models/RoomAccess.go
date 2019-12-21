package models

import (
	"github.com/jinzhu/gorm"
)

type RoomAccess struct {
	gorm.Model
	ID         uint   `json:"id"`
	AccessType string `json:"access_type"`
	Camera     Camera `json:"camera"`
}
