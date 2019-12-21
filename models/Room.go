package models

import (
	"github.com/jinzhu/gorm"
)

type Room struct {
	gorm.Model
	ID     uint   `json:"id"`
	User   User   `json:"user"`
	Camera Camera `json:"camera"`
}
