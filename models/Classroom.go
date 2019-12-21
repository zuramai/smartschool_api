package models

import (
	"github.com/jinzhu/gorm"
)

type Classroom struct {
	gorm.Model
	ID       uint   `json:"id" gorm:"primary_key"`
	Name     string `json:"name" gorm:"primary_key"`
	Location int    `json:"user_id"`
	Note     string `json:"note"`
	Address  string `json:"address"`
}

func (Classroom) TableName() string {
	return "classroom"
}
