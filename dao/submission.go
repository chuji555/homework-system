package dao

import (
	"github.com/chuji555/homework-system/models"
	"gorm.io/gorm"
)

// 创建作业提交记录
func CreateSubmission(submission *models.Submission) error {
	return DB.Create(submission).Error
}

// 根据学生ID+作业ID查询提交记录（防止重复提交）
func GetSubmissionByStudentAndHomework(studentID, homeworkID int64) (*models.Submission, error) {
	var sub models.Submission
	err := DB.Where("student_id = ? AND homework_id = ?", studentID, homeworkID).First(&sub).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &sub, err
}

// 根据学生ID分页查询提交记录
func ListSubmissionByStudentID(studentID int64, page, pageSize int) ([]models.Submission, int64, error) {
	var list []models.Submission
	var total int64

	// 先查总数
	if err := DB.Model(&models.Submission{}).Where("student_id = ?", studentID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := DB.Preload("Homework"). // 关联查询作业信息
					Where("student_id = ?", studentID).
					Order("submitted_at DESC").
					Limit(pageSize).
					Offset(offset).
					Find(&list).Error

	return list, total, err
}

// 根据作业ID分页查询提交记录（管理员用）
func ListSubmissionByHomeworkID(homeworkID int64, page, pageSize int) ([]models.Submission, int64, error) {
	var list []models.Submission
	var total int64

	if err := DB.Model(&models.Submission{}).Where("homework_id = ?", homeworkID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := DB.Preload("Student"). // 关联查询学生信息
					Where("homework_id = ?", homeworkID).
					Order("submitted_at DESC").
					Limit(pageSize).
					Offset(offset).
					Find(&list).Error

	return list, total, err
}

// 根据ID查询提交记录
func GetSubmissionByID(subID int64) (*models.Submission, error) {
	var sub models.Submission
	err := DB.First(&sub, subID).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &sub, err
}

// 更新提交记录（批改作业/标记优秀）
func UpdateSubmission(submission *models.Submission) error {
	return DB.Save(submission).Error
}

// 查询优秀作业（所有学生可见）
func ListExcellentSubmission(page, pageSize int) ([]models.Submission, int64, error) {
	var list []models.Submission
	var total int64

	if err := DB.Model(&models.Submission{}).Where("is_excellent = ?", true).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := DB.Preload("Homework").
		Preload("Student").
		Where("is_excellent = ?", true).
		Order("submitted_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&list).Error

	return list, total, err
}
