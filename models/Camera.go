package models

import "github.com/jinzhu/gorm"

type Camera struct {
	gorm.Model
	Name     string `json:"name" `
	Location string `json:"location" gorm:"type:varchar(100)"`
	Note     string `json:"note" gorm:"type:text"`
	Status   bool   `json:"status" `
}

func (Camera) TableName() string {
	return "cameras"
}
