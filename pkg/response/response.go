package response

import (
	"github.com/chuji555/homework-system/pkg/errcode"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code    errcode.ErrCode `json:"code"`
	Message string          `json:"message"`
	Data    interface{}     `json:"data"`
}

// 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    errcode.Success,
		Message: errcode.Success.Msg(),
		Data:    data,
	})
}

// 错误响应
func Error(c *gin.Context, code errcode.ErrCode) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: code.Msg(),
		Data:    nil,
	})
}

// 分页响应（作业列表、提交列表用）
type PageResponse struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}
