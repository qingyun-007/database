package router

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "getcharzp.cn/docs"
	"getcharzp.cn/middlewares"
	"getcharzp.cn/service"
)

func Router() *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.Cors())

	// Swagger 配置
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	rs := r.Group("/api")
	// 路由规则
	// 公有方法
	// 处理器信息
	rs.GET("/cpuInfo", service.GetCpuInfo)

	// 问题
	rs.GET("/problem-list", service.GetProblemList)
	rs.GET("/problem-detail", service.GetProblemDetail) // 根据id查询题目
	// 用户
	rs.GET("/user-detail", service.GetUserDetail)
	rs.POST("/login", service.Login)
	rs.POST("/send-code", service.SendCode)
	rs.POST("/register", service.Register)

	// 排行榜
	rs.GET("/rank-list", service.GetRankList) // 就是简单的数据库查询
	// 提交记录
	rs.GET("/submit-list", service.GetSubmitList)
	// 分类列表
	rs.GET("/category-list", service.GetCategoryList)

	// 管理员私有方法
	authAdmin := rs.Group("/admin", middlewares.AuthAdminCheck())
	//authAdmin := r.Group("/admin")
	// 问题创建
	authAdmin.POST("/problem-create", service.ProblemCreate)
	// 问题修改
	authAdmin.PUT("/problem-modify", service.ProblemModify)
	// 分类创建
	authAdmin.POST("/category-create", service.CategoryCreate)
	// 分类修改
	authAdmin.PUT("/category-modify", service.CategoryModify)
	// 分类删除
	authAdmin.DELETE("/category-delete", service.CategoryDelete)
	// 获取测试案例
	authAdmin.GET("/test-case", service.GetTestCase)

	// 用户私有方法
	authUser := rs.Group("/user", middlewares.AuthUserCheck())
	// 代码提交
	authUser.POST("/submit", service.Submit)
	return r
}
