package models

import (
	"time"
)

type Homework struct {
	ID          int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string     `gorm:"size:200;not null" json:"title"`
	Description string     `gorm:"type:text;not null" json:"description"`
	Department  Department `gorm:"type:enum('backend','frontend','sre','product','design','android','ios');not null" json:"department"`
	CreatorID   int64      `gorm:"not null" json:"creator_id"`
	Deadline    time.Time  `gorm:"not null" json:"deadline"`
	AllowLate   bool       `gorm:"default:false" json:"allow_late"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	// 关联发布者（后续查询用）
	Creator User `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
}

func (h *Homework) DepartmentLabel() string {
	switch h.Department {
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
