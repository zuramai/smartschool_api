package models

import "time"

type Camera struct {
	ID        uint       `json:"id" gorm:"type:bigint(20) unsigned auto_increment;not null;primary_key"`
	Name      string     `json:"name" gorm:"primary_key"`
	Location  string     `json:"location" gorm:"type:varchar(100)"`
	Note      string     `json:"note" gorm:"type:text"`
	Status    bool       `json:"status" `
	CreatedAt time.Time  `json:"created_at" gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func (Camera) TableName() string {
	return "cameras"
}
