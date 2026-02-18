package models

import (
	"time"

	"gorm.io/gorm"
)

type Department string

const (
	Backend  Department = "backend"
	Frontend Department = "frontend"
	SRE      Department = "sre"
	Product  Department = "product"
	Design   Department = "design"
	Android  Department = "android"
	iOS      Department = "ios"
)

type Role string

const (
	Student Role = "student"
	Admin   Role = "admin"
)

type User struct {
	ID         int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username   string     `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Password   string     `gorm:"size:255;not null" json:"-"` // 序列化时隐藏密码
	Nickname   string     `gorm:"size:50;not null" json:"nickname"`
	Role       Role       `gorm:"type:enum('student','admin');not null" json:"role"`
	Department Department `gorm:"type:enum('backend','frontend','sre','product','design','android','ios');not null" json:"department"`
	Email      string     `gorm:"size:100" json:"email"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	// 软删除标记
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (u *User) DepartmentLabel() string {
	switch u.Department {
	case Backend:
		return "后端"
	case Frontend:
		return "前端"
	case SRE:
		return "SRE"
	case Product:
		return "产品"
	case Design:
		return "视觉设计"
	case Android:
		return "Android"
	case iOS:
		return "iOS"
	default:
		return ""
	}
}
