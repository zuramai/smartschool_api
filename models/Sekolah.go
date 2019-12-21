package models

type Sekolah struct {
	ID   int64  `json:"id" gorm:"type:bigint(20) unsigned auto_increment;not null;primary_key"`
	Nama string `json:"nama" gorm:"column:nama"`
}

func (Sekolah) TableName() string {
	return "sekolah"
}
