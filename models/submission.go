package models

import (
	"time"
)

type Submission struct {
	ID          int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	HomeworkID  int64      `gorm:"not null;index" json:"homework_id"`
	StudentID   int64      `gorm:"not null;index" json:"student_id"`
	Content     string     `gorm:"type:text;not null" json:"content"`
	FileURL     string     `gorm:"size:500" json:"file_url"`
	IsLate      bool       `gorm:"default:false" json:"is_late"`
	Score       *int       `json:"score,omitempty"` // 分数可选（批改后才有）
	Comment     string     `gorm:"type:text" json:"comment,omitempty"`
	IsExcellent bool       `gorm:"default:false" json:"is_excellent"`
	ReviewerID  *int64     `json:"reviewer_id,omitempty"`
	SubmittedAt time.Time  `json:"submitted_at"`
	ReviewedAt  *time.Time `json:"reviewed_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	// 关联作业和学生
	Homework Homework `gorm:"foreignKey:HomeworkID" json:"homework,omitempty"`
	Student  User     `gorm:"foreignKey:StudentID" json:"student,omitempty"`
}
