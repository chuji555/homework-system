package handler

import (
	"github.com/chuji555/homework-system/models"
	"github.com/chuji555/homework-system/pkg/errcode"
	"github.com/chuji555/homework-system/pkg/response"
	"github.com/chuji555/homework-system/service"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

// CreateHomeworkRequest 管理员创建作业的请求参数
type CreateHomeworkRequest struct {
	Title       string    `json:"title" binding:"required,max=200"`                                                    // 作业标题（必填，最长200字符）
	Description string    `json:"description" binding:"required"`                                                      // 作业描述（必填）
	Department  string    `json:"department" binding:"required,oneof=backend frontend sre product design android ios"` // 所属部门（必填，限定枚举值）
	Deadline    time.Time `json:"deadline" binding:"required"`                                                         // 截止时间（必填）
	AllowLate   bool      `json:"allow_late" binding:"omitempty"`                                                      // 是否允许迟交（可选，默认false）
}

// UpdateHomeworkRequest 管理员修改作业的请求参数（和创建类似，字段可选）
type UpdateHomeworkRequest struct {
	Title       string     `json:"title" binding:"omitempty,max=200"`
	Description string     `json:"description" binding:"omitempty"`
	Department  string     `json:"department" binding:"omitempty,oneof=backend frontend sre product design android ios"`
	Deadline    *time.Time `json:"deadline" binding:"omitempty"` // 用指针，区分“不传”和“传空”
	AllowLate   *bool      `json:"allow_late" binding:"omitempty"`
}

// -------------------------- 核心接口实现 --------------------------

// CreateHomework 管理员创建作业
func CreateHomework(c *gin.Context) {
	// 1. 获取当前登录的管理员ID（从上下文取，AuthMiddleware已存入）
	creatorID, exists := c.Get("userID")
	if !exists {
		response.Error(c, errcode.AuthError)
		return
	}

	// 2. 绑定并校验请求参数（binding标签会自动校验必填/枚举/长度）
	var req CreateHomeworkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 参数校验失败，返回参数错误码
		response.Error(c, errcode.ParamError)
		return
	}

	// 3. 调用service层的创建作业逻辑
	errCode := service.CreateHomework(
		req.Title,
		req.Description,
		req.Department,
		creatorID.(int64), // 类型断言：上下文存的是interface{}，转成int64
		req.Deadline,
		req.AllowLate,
	)

	// 4. 根据业务逻辑结果返回响应
	if errCode != errcode.Success {
		response.Error(c, errCode)
		return
	}

	// 创建成功，返回提示
	response.Success(c, gin.H{"msg": "作业创建成功"})
}

// UpdateHomework 管理员修改作业
func UpdateHomework(c *gin.Context) {
	// 1. 获取路径参数：作业ID
	homeworkIDStr := c.Param("id")
	homeworkID, err := strconv.ParseInt(homeworkIDStr, 10, 64)
	if err != nil || homeworkID <= 0 {
		// ID格式错误，返回参数错误
		response.Error(c, errcode.ParamError)
		return
	}

	// 2. 绑定并校验修改参数
	var req UpdateHomeworkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ParamError)
		return
	}

	// 3. 调用service层修改逻辑
	errCode := service.UpdateHomework(
		homeworkID,
		req.Title,
		req.Description,
		req.Department,
		req.Deadline,
		req.AllowLate,
	)

	// 4. 返回响应
	if errCode != errcode.Success {
		response.Error(c, errCode)
		return
	}
	response.Success(c, gin.H{"msg": "作业修改成功"})
}

// DeleteHomework 管理员删除作业（软删除）
func DeleteHomework(c *gin.Context) {
	// 1. 获取作业ID
	homeworkIDStr := c.Param("id")
	homeworkID, err := strconv.ParseInt(homeworkIDStr, 10, 64)
	if err != nil || homeworkID <= 0 {
		response.Error(c, errcode.ParamError)
		return
	}

	// 2. 调用service层删除逻辑
	errCode := service.DeleteHomework(homeworkID)

	// 3. 返回响应
	if errCode != errcode.Success {
		response.Error(c, errCode)
		return
	}
	response.Success(c, gin.H{"msg": "作业删除成功"})
}

// ListHomework 所有人查询作业列表（支持按部门筛选+分页）
func ListHomework(c *gin.Context) {
	// 1. 获取分页参数（默认第1页，每页10条）
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	// 分页参数校验（防止负数/过大）
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	// 2. 获取筛选参数：部门（可选，不传则查所有）
	department := c.Query("department")
	// 校验部门参数（如果传了，必须是枚举值）
	if department != "" {
		validDept := false
		// 遍历部门枚举，校验合法性
		for _, d := range []string{"backend", "frontend", "sre", "product", "design", "android", "ios"} {
			if department == d {
				validDept = true
				break
			}
		}
		if !validDept {
			response.Error(c, errcode.ParamError)
			return
		}
	}

	// 3. 调用service层查询列表逻辑
	list, total, errCode := service.ListHomework(department, page, pageSize)
	if errCode != errcode.Success {
		response.Error(c, errCode)
		return
	}

	// 4. 构造分页响应（统一格式）
	resp := response.PageResponse{
		List:     formatHomeworkList(list), // 格式化列表（补充部门中文标签）
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}
	response.Success(c, resp)
}

// GetHomework 所有人查询作业详情
func GetHomework(c *gin.Context) {
	// 1. 获取作业ID
	homeworkIDStr := c.Param("id")
	homeworkID, err := strconv.ParseInt(homeworkIDStr, 10, 64)
	if err != nil || homeworkID <= 0 {
		response.Error(c, errcode.ParamError)
		return
	}

	// 2. 调用service层查询详情逻辑
	homework, errCode := service.GetHomeworkByID(homeworkID)
	if errCode != errcode.Success {
		response.Error(c, errCode)
		return
	}

	// 3. 格式化响应（补充部门中文标签）
	resp := gin.H{
		"id":               homework.ID,
		"title":            homework.Title,
		"description":      homework.Description,
		"department":       homework.Department,
		"department_label": homework.DepartmentLabel(), // 部门中文标签
		"creator_id":       homework.CreatorID,
		"creator_nickname": homework.Creator.Nickname, // 发布者昵称（关联查询）
		"deadline":         homework.Deadline,
		"allow_late":       homework.AllowLate,
		"created_at":       homework.CreatedAt,
		"updated_at":       homework.UpdatedAt,
	}

	response.Success(c, resp)
}

// formatHomeworkList 格式化作业列表，补充部门中文标签
func formatHomeworkList(homeworks []models.Homework) []gin.H {
	var list []gin.H
	for _, h := range homeworks {
		list = append(list, gin.H{
			"id":               h.ID,
			"title":            h.Title,
			"description":      h.Description,
			"department":       h.Department,
			"department_label": h.DepartmentLabel(),
			"creator_id":       h.CreatorID,
			"deadline":         h.Deadline,
			"allow_late":       h.AllowLate,
			"created_at":       h.CreatedAt,
		})
	}
	return list
}
