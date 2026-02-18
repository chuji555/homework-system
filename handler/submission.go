package handler

import (
	"github.com/chuji555/homework-system/pkg/errcode"
	"github.com/chuji555/homework-system/pkg/response"
	"github.com/chuji555/homework-system/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

// 提交作业的请求参数
type CreateSubmissionRequest struct {
	HomeworkID int64  `json:"homework_id" binding:"required"` // 作业ID
	Content    string `json:"content" binding:"required"`     // 提交内容
	FileURL    string `json:"file_url"`                       // 附件URL（可选）
}

// 学生提交作业
func CreateSubmission(c *gin.Context) {
	// 1. 获取当前登录学生ID
	studentID, _ := c.Get("userID")
	if studentID == nil {
		response.Error(c, errcode.AuthError)
		return
	}

	// 2. 校验请求参数
	var req CreateSubmissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ParamError)
		return
	}

	// 3. 调用业务逻辑
	errCode := service.CreateSubmission(
		studentID.(int64),
		req.HomeworkID,
		req.Content,
		req.FileURL,
	)
	if errCode != errcode.Success {
		response.Error(c, errCode)
		return
	}

	// 4. 返回成功响应
	response.Success(c, gin.H{"msg": "提交成功"})
}

// 学生查询我的提交记录
func ListMySubmission(c *gin.Context) {
	// 1. 获取当前学生ID
	studentID, _ := c.Get("userID")
	if studentID == nil {
		response.Error(c, errcode.AuthError)
		return
	}

	// 2. 获取分页参数（默认第1页，每页10条）
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	// 3. 调用业务逻辑
	list, total, errCode := service.ListMySubmission(studentID.(int64), page, pageSize)
	if errCode != errcode.Success {
		response.Error(c, errCode)
		return
	}

	// 4. 构造分页响应
	resp := response.PageResponse{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}
	response.Success(c, resp)
}

// 管理员根据作业ID查询提交记录
func ListSubmissionByHomework(c *gin.Context) {
	// 1. 获取作业ID
	homeworkIDStr := c.Param("homework_id")
	homeworkID, err := strconv.ParseInt(homeworkIDStr, 10, 64)
	if err != nil {
		response.Error(c, errcode.ParamError)
		return
	}

	// 2. 分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	// 3. 调用业务逻辑
	list, total, errCode := service.ListSubmissionByHomework(homeworkID, page, pageSize)
	if errCode != errcode.Success {
		response.Error(c, errCode)
		return
	}

	resp := response.PageResponse{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}
	response.Success(c, resp)
}

// 管理员批改作业
type ReviewSubmissionRequest struct {
	Score   int    `json:"score" binding:"required,min=0,max=100"` // 分数0-100
	Comment string `json:"comment"`                                // 批改评语
}

func ReviewSubmission(c *gin.Context) {
	// 1. 获取提交ID
	subIDStr := c.Param("id")
	subID, err := strconv.ParseInt(subIDStr, 10, 64)
	if err != nil {
		response.Error(c, errcode.ParamError)
		return
	}

	// 2. 获取当前管理员ID
	reviewerID, _ := c.Get("userID")
	if reviewerID == nil {
		response.Error(c, errcode.AuthError)
		return
	}

	// 3. 校验参数
	var req ReviewSubmissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ParamError)
		return
	}

	// 4. 调用业务逻辑
	errCode := service.ReviewSubmission(subID, req.Score, req.Comment, reviewerID.(int64))
	if errCode != errcode.Success {
		response.Error(c, errCode)
		return
	}

	response.Success(c, gin.H{"msg": "批改成功"})
}

// 管理员标记优秀作业
type MarkExcellentRequest struct {
	IsExcellent bool `json:"is_excellent" binding:"required"`
}

func MarkExcellent(c *gin.Context) {
	// 1. 获取提交ID
	subIDStr := c.Param("id")
	subID, err := strconv.ParseInt(subIDStr, 10, 64)
	if err != nil {
		response.Error(c, errcode.ParamError)
		return
	}

	// 2. 校验参数
	var req MarkExcellentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ParamError)
		return
	}

	// 3. 调用业务逻辑
	errCode := service.MarkExcellent(subID, req.IsExcellent)
	if errCode != errcode.Success {
		response.Error(c, errCode)
		return
	}

	response.Success(c, gin.H{"msg": "标记成功"})
}

// 查询优秀作业（所有人可见）
func ListExcellentSubmission(c *gin.Context) {
	// 1. 分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	// 2. 调用业务逻辑
	list, total, errCode := service.ListExcellentSubmission(page, pageSize)
	if errCode != errcode.Success {
		response.Error(c, errCode)
		return
	}

	resp := response.PageResponse{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}
	response.Success(c, resp)
}
