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
	ID        uint       `json:"id" gorm:"type:bigint(20) unsigned auto_increment;not null;primary_key"`
	Username  string     `json:"username" gorm:"type:varchar(100);unique_index"`
	Password  string     `json:"-" gorm:"type:varchar(100);"`
	Name      string     `json:"name" gorm:"size:255"`
	Phone     string     `json:"phone" gorm:"size:14"`
	Photo     string     `json:"-" gorm:"type:varchar(100)"`
	Role      string     `json:"role" gorm:"type:enum('Superadmin','Moderator','Teacher','Student');unique_index"`
	Status    bool       `json:"status" gorm:"type:boolean"`
	Token     string     `json:"token";sql:"-"`
	CreatedAt time.Time  `json:"createdAt" gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time  `json:"updatedAt" gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt *time.Time `json:"deleted_at"`
}
