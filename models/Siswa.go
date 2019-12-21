package models

import "time"

type Siswa struct {
	ID              int64     `json:"id" gorm:"type:bigint(20) unsigned auto_increment;not null;primary_key"`
	UserID          int64     `json:"user_id" gorm:"type:bigint;index;foreignkey"`
	NoInduk         int64     `json:"no_induk"`
	NIS             string    `json:"nis" gorm:"type:varchar(100);unique_index"`
	NISN            string    `json:"nisn" gorm:"type:varchar(100);unique_index"`
	NIK             string    `json:"nik" gorm:"size:100;unique_index"`
	NamaLengkap     string    `json:"nama_lengkap" gorm:"size:14"`
	Kelamin         string    `json:"kelamin" gorm:"type:enum('L','P')"`
	TempatLahir     string    `json:"tempat_lahir" gorm:"type:varchar(100)"`
	TglLahir        bool      `json:"tgl_lahir" gorm:"type:date"`
	Agama           bool      `json:"agama" gorm:"type:varchar(100)"`
	Kewarganegaraan bool      `json:"kewarganegaraan" gorm:"type:varchar(100)"`
	Status          bool      `json:"status" gorm:"type:boolean"`
	CreatedAt       time.Time `json:"createdAt" gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time `json:"updatedAt" gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt       *time.Time
}

func (Siswa) TableName() string {
	return "siswa"
}
