package handler

import (
	"github.com/chuji555/homework-system/pkg/errcode"
	"github.com/chuji555/homework-system/pkg/response"
	"github.com/chuji555/homework-system/service"

	"github.com/gin-gonic/gin"
)

// 注册请求参数
type RegisterRequest struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Nickname   string `json:"nickname" binding:"required"`
	Department string `json:"department" binding:"required"`
}

// 注册接口
func Register(c *gin.Context) {
	var req RegisterRequest
	// 参数校验（binding:"required"会自动校验）
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ParamError)
		return
	}
	// 调用业务逻辑
	user, errCode := service.Register(req.Username, req.Password, req.Nickname, req.Department)
	if errCode != errcode.Success {
		response.Error(c, errCode)
		return
	}
	// 构造响应（包含department_label）
	resp := gin.H{
		"id":               user.ID,
		"username":         user.Username,
		"nickname":         user.Nickname,
		"role":             user.Role,
		"department":       user.Department,
		"department_label": user.DepartmentLabel(),
	}
	response.Success(c, resp)
}

// 登录请求参数
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 登录接口
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ParamError)
		return
	}
	accessToken, refreshToken, user, errCode := service.Login(req.Username, req.Password)
	if errCode != errcode.Success {
		response.Error(c, errCode)
		return
	}
	resp := gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user": gin.H{
			"id":               user.ID,
			"username":         user.Username,
			"nickname":         user.Nickname,
			"role":             user.Role,
			"department":       user.Department,
			"department_label": user.DepartmentLabel(),
		},
	}
	response.Success(c, resp)
}

// 刷新Token接口
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func RefreshToken(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ParamError)
		return
	}
	newAccessToken, newRefreshToken, errCode := service.RefreshToken(req.RefreshToken)
	if errCode != errcode.Success {
		response.Error(c, errCode)
		return
	}
	resp := gin.H{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	}
	response.Success(c, resp)
}

// 注销账号接口
type LogoutRequest struct {
	Password string `json:"password" binding:"required"`
}

func Logout(c *gin.Context) {
	// 获取当前用户ID（AuthMiddleware存入的）
	userID, _ := c.Get("userID")
	var req LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ParamError)
		return
	}
	errCode := service.Logout(userID.(int64))
	if errCode != errcode.Success {
		response.Error(c, errCode)
		return
	}
	response.Success(c, nil)
}

// 获取用户信息接口
func GetProfile(c *gin.Context) {
	userID, _ := c.Get("userID")
	username, _ := c.Get("username")
	role, _ := c.Get("role")
	department, _ := c.Get("department")
	// 查询用户详情（获取邮箱等）
	user, err := service.GetUserByID(userID.(int64))
	if err != nil || user == nil {
		response.Error(c, errcode.DataNotFound)
		return
	}
	resp := gin.H{
		"id":               user.ID,
		"username":         username,
		"nickname":         user.Nickname,
		"role":             role,
		"department":       department,
		"department_label": user.DepartmentLabel(),
		"email":            user.Email,
	}
	response.Success(c, resp)
}
