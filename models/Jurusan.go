package models

import (
	"github.com/jinzhu/gorm"
)

type Jurusan struct {
	gorm.Model
	Singkatan string
	Kode      string
	Paket     string
}
