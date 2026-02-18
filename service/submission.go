package service

import (
	"time"

	"github.com/chuji555/homework-system/dao"
	"github.com/chuji555/homework-system/models"
	"github.com/chuji555/homework-system/pkg/errcode"
)

// 提交作业的业务逻辑
func CreateSubmission(studentID, homeworkID int64, content, fileURL string) errcode.ErrCode {
	// 1. 检查是否已提交
	sub, err := dao.GetSubmissionByStudentAndHomework(studentID, homeworkID)
	if err != nil {
		return errcode.DBError
	}
	if sub != nil {
		return errcode.ParamError // 自定义"已提交过该作业"的错误码也可以
	}

	// 2. 检查作业截止时间（简化版：实际应先查询作业信息）
	// 这里先跳过作业查询，直接标记是否迟交（后续可以补充）
	isLate := false // 实际逻辑：对比提交时间和作业deadline

	// 3. 创建提交记录
	submission := &models.Submission{
		HomeworkID:  homeworkID,
		StudentID:   studentID,
		Content:     content,
		FileURL:     fileURL,
		IsLate:      isLate,
		SubmittedAt: time.Now(),
	}

	if err := dao.CreateSubmission(submission); err != nil {
		return errcode.DBError
	}

	return errcode.Success
}

// 查询我的提交记录
func ListMySubmission(studentID int64, page, pageSize int) ([]models.Submission, int64, errcode.ErrCode) {
	list, total, err := dao.ListSubmissionByStudentID(studentID, page, pageSize)
	if err != nil {
		return nil, 0, errcode.DBError
	}
	return list, total, errcode.Success
}

// 管理员查询作业的所有提交
func ListSubmissionByHomework(homeworkID int64, page, pageSize int) ([]models.Submission, int64, errcode.ErrCode) {
	list, total, err := dao.ListSubmissionByHomeworkID(homeworkID, page, pageSize)
	if err != nil {
		return nil, 0, errcode.DBError
	}
	return list, total, errcode.Success
}

// 批改作业
func ReviewSubmission(subID int64, score int, comment string, reviewerID int64) errcode.ErrCode {
	// 1. 查询提交记录
	sub, err := dao.GetSubmissionByID(subID)
	if err != nil {
		return errcode.DBError
	}
	if sub == nil {
		return errcode.DataNotFound
	}

	// 2. 更新批改信息
	now := time.Now()
	sub.Score = &score
	sub.Comment = comment
	sub.ReviewerID = &reviewerID
	sub.ReviewedAt = &now

	if err := dao.UpdateSubmission(sub); err != nil {
		return errcode.DBError
	}

	return errcode.Success
}

// 标记优秀作业
func MarkExcellent(subID int64, isExcellent bool) errcode.ErrCode {
	sub, err := dao.GetSubmissionByID(subID)
	if err != nil {
		return errcode.DBError
	}
	if sub == nil {
		return errcode.DataNotFound
	}

	sub.IsExcellent = isExcellent
	if err := dao.UpdateSubmission(sub); err != nil {
		return errcode.DBError
	}

	return errcode.Success
}

// 查询优秀作业
func ListExcellentSubmission(page, pageSize int) ([]models.Submission, int64, errcode.ErrCode) {
	list, total, err := dao.ListExcellentSubmission(page, pageSize)
	if err != nil {
		return nil, 0, errcode.DBError
	}
	return list, total, errcode.Success
}
