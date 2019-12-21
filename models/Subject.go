package models

type Subject struct {
	ID   int64  `json:"id" gorm:"type:bigint(20) unsigned auto_increment;not null;primary_key"`
	Code string `json:"code"`
	Name string `json:"name"`
}
