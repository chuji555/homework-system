package dao

import (
	"github.com/chuji555/homework-system/models"
	"gorm.io/gorm"
)

// CreateHomework 创建作业
func CreateHomework(homework *models.Homework) error {
	return DB.Create(homework).Error
}

// UpdateHomework 修改作业
func UpdateHomework(homework *models.Homework) error {
	return DB.Save(homework).Error
}

// DeleteHomework 软删除作业
func DeleteHomework(homeworkID int64) error {
	return DB.Delete(&models.Homework{}, homeworkID).Error
}

// ListHomework 分页查询作业（支持部门筛选）
func ListHomework(department string, page, pageSize int) ([]models.Homework, int64, error) {
	var list []models.Homework
	var total int64

	// 构建查询条件
	query := DB.Model(&models.Homework{})
	if department != "" {
		query = query.Where("department = ?", department)
	}

	// 先查总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&list).Error

	return list, total, err
}

// GetHomeworkByID 根据ID查询作业（关联查询发布者信息）
func GetHomeworkByID(homeworkID int64) (*models.Homework, error) {
	var homework models.Homework
	err := DB.Preload("Creator"). // 关联查询发布者信息（User表）
					Where("id = ?", homeworkID).
					First(&homework).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &homework, err
}
