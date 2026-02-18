package service

import (
	"github.com/chuji555/homework-system/dao"
	"github.com/chuji555/homework-system/models"
	"github.com/chuji555/homework-system/pkg/errcode"
	"time"
)

// CreateHomework 创建作业
func CreateHomework(title, desc, dept string, creatorID int64, deadline time.Time, allowLate bool) errcode.ErrCode {
	// 先声明并初始化 homework 变量
	homework := &models.Homework{
		Title:       title,
		Description: desc,
		Department:  models.Department(dept),
		CreatorID:   creatorID,
		Deadline:    deadline,
		AllowLate:   allowLate,
	}

	// 调用 dao 层创建
	if err := dao.CreateHomework(homework); err != nil {
		return errcode.DBError
	}
	return errcode.Success
}

// UpdateHomework 修改作业
func UpdateHomework(homeworkID int64, title, desc, dept string, deadline *time.Time, allowLate *bool) errcode.ErrCode {
	// 1. 先查询作业是否存在
	homework, err := dao.GetHomeworkByID(homeworkID)
	if err != nil {
		return errcode.DBError
	}
	if homework == nil {
		return errcode.DataNotFound
	}

	// 2. 只更新传了的字段（指针判断是否传值）
	if title != "" {
		homework.Title = title
	}
	if desc != "" {
		homework.Description = desc
	}
	if dept != "" {
		homework.Department = models.Department(dept)
	}
	if deadline != nil {
		homework.Deadline = *deadline
	}
	if allowLate != nil {
		homework.AllowLate = *allowLate
	}

	// 3. 调用 dao 层修改
	if err := dao.UpdateHomework(homework); err != nil {
		return errcode.DBError
	}
	return errcode.Success
}

// DeleteHomework 删除作业
func DeleteHomework(homeworkID int64) errcode.ErrCode {
	// 先检查作业是否存在
	homework, err := dao.GetHomeworkByID(homeworkID)
	if err != nil {
		return errcode.DBError
	}
	if homework == nil {
		return errcode.DataNotFound
	}

	// 调用 dao 层删除
	if err := dao.DeleteHomework(homeworkID); err != nil {
		return errcode.DBError
	}
	return errcode.Success
}

// ListHomework 分页查询作业列表
func ListHomework(department string, page, pageSize int) ([]models.Homework, int64, errcode.ErrCode) {
	list, total, err := dao.ListHomework(department, page, pageSize)
	if err != nil {
		return nil, 0, errcode.DBError
	}
	return list, total, errcode.Success
}

// GetHomeworkByID 查询作业详情
func GetHomeworkByID(homeworkID int64) (*models.Homework, errcode.ErrCode) {
	homework, err := dao.GetHomeworkByID(homeworkID)
	if err != nil {
		return nil, errcode.DBError
	}
	if homework == nil {
		return nil, errcode.DataNotFound
	}
	return homework, errcode.Success
}
