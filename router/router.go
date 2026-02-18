package router

import (
	"github.com/chuji555/homework-system/handler"
	"github.com/chuji555/homework-system/middleware"
	"github.com/gin-gonic/gin"
)

// 初始化路由
func InitRouter() *gin.Engine {
	r := gin.Default()
	// 公开接口（无需认证）
	publicGroup := r.Group("/")
	{
		publicGroup.POST("/user/register", handler.Register)
		publicGroup.POST("/user/login", handler.Login)
		publicGroup.POST("/user/refresh", handler.RefreshToken)
	}

	// 需要认证的接口（所有请求都要带AccessToken）
	authGroup := r.Group("/")
	authGroup.Use(middleware.AuthMiddleware())
	{
		// 用户模块
		userGroup := authGroup.Group("/user")
		{
			userGroup.GET("/profile", handler.GetProfile)
			userGroup.DELETE("/account", handler.Logout)
		}
		// 作业模块
		homeworkGroup := authGroup.Group("/homework")
		{
			// 老登才能发布/修改/删除
			homeworkGroup.POST("", middleware.AdminMiddleware(), handler.CreateHomework)
			homeworkGroup.PUT("/:id", middleware.AdminMiddleware(), handler.UpdateHomework)
			homeworkGroup.DELETE("/:id", middleware.AdminMiddleware(), handler.DeleteHomework)
			// 所有人都能查列表和详情
			homeworkGroup.GET("", handler.ListHomework)
			homeworkGroup.GET("/:id", handler.GetHomework)
		}
		// 提交模块
		submissionGroup := authGroup.Group("/submission")
		{
			// 小登才能提交
			submissionGroup.POST("", middleware.StudentMiddleware(), handler.CreateSubmission)
			// 小登查自己的提交
			submissionGroup.GET("/my", middleware.StudentMiddleware(), handler.ListMySubmission)
			// 老登查部门提交、批改、标记优秀
			submissionGroup.GET("/homework/:homework_id", middleware.AdminMiddleware(), handler.ListSubmissionByHomework)
			submissionGroup.PUT("/:id/review", middleware.AdminMiddleware(), handler.ReviewSubmission)
			submissionGroup.PUT("/:id/excellent", middleware.AdminMiddleware(), handler.MarkExcellent)
			// 所有人查优秀作业
			submissionGroup.GET("/excellent", handler.ListExcellentSubmission)
		}
	}
	return r
}
