package models

import (
	"time"
)

// Attendance : Base struct table attendance
type Attendance struct {
	UserID       string    `json:"-" bson:"user_id"`
	PictureTaken string    `json:"picture_taken" bson:"picture_taken"`
	AttendAt     time.Time `json:"attend_at" bson:"attend_at"`
	Keterangan   string    `json:"status"`
}

// Attendance : attendance for json response
type AttendanceJSON struct {
	User      UserJSON  `json:"user"`
	PhotoName string    `json:"photo_name" gorm:"type:varchar(255)"`
	AttendAt  time.Time `json:"attend_at" gorm:"type:timestamp" sql:"DEFAULT:current_timestamp"`
	Status    string    `json:"status"`
}

// TableName : Set the name of table Attendance
func (Attendance) TableName() string {
	return "attendances"
}
