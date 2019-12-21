package models

import (
	"time"
)

type Attendance struct {
	ID        int64     `json:"id" gorm:"type:bigint(20) unsigned auto_increment;not null;primary_key"`
	UserID    int64     `json:"user_id" gorm:"type:bigint;index;foreign_key"`
	PhotoName string    `json:"photo_name" gorm:"type:varchar(255)"`
	AttendAt  time.Time `json:"attend_at" gorm:"type:datetime" sql:"DEFAULT:current_timestamp"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

func (Attendance) TableName() string {
	return "attendances"
}
