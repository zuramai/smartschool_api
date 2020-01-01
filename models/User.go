package models

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Token struct {
	UserID uint
	jwt.StandardClaims
}

type User struct {
	ID        uint       `json:"id" gorm:"primary_key"`
	Username  string     `json:"username" gorm:"size:100;unique_index"`
	Password  string     `json:"-" gorm:"size:100;"`
	Name      string     `json:"name" gorm:"size:255"`
	Phone     string     `json:"phone" gorm:"size:14"`
	PhotoName string     `json:"photo_name" gorm:"size:100"`
	RoleID    uint       `json:"-" `
	Role      Role       `json:"role" gorm:"foreignkey:RoleID"`
	Status    bool       `json:"status" gorm:"type:boolean"`
	Token     string     `json:"token" sql:"-"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type UserJSON struct {
	ID        uint   `json:"id" gorm:"primary_key"`
	Username  string `json:"username" gorm:"size:100;unique_index"`
	Name      string `json:"name" gorm:"size:255"`
	PhotoName string `json:"photo_name" gorm:"size:100"`
}
