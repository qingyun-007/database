package router

import (
	_ "getcharzp.cn/docs"
)

//func Router1() *gin.Engine {
//	r := gin.Default()
//	r.Use(middlewares.Cors())
//
//	// Swagger 配置
//	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
//
//	rs := r.Group("/api")
//	// 路由规则
//	// 公有方法
//	// 处理器信息
//	rs.GET("/cpuInfo", service.GetCpuInfo)
//
//	// 问题
//	rs.GET("/problem-list", service.GetProblemList)
//	rs.GET("/problem-detail", service.GetProblemDetail) // 根据id查询题目
//	authAdmin.POST("/problem-create", service.ProblemCreate)
//	authAdmin.PUT("/problem-modify", service.ProblemModify)
//
//	// 分类列表
//	rs.GET("/category-list", service.GetCategoryList)
//	// 分类创建
//	authAdmin.POST("/category-create", service.CategoryCreate)
//	// 分类修改
//	authAdmin.PUT("/category-modify", service.CategoryModify)
//	// 分类删除
//	authAdmin.DELETE("/category-delete", service.CategoryDelete)
//
//	// 用户
//	rs.GET("/user-detail", service.GetUserDetail)
//	rs.POST("/login", service.Login)
//	rs.POST("/send-code", service.SendCode)
//	rs.POST("/register", service.Register)
//
//	// 排行榜
//	rs.GET("/rank-list", service.GetRankList) // 就是简单的数据库查询
//
//	// 代码提交
//	authUser.POST("/submit", service.Submit)
//	// 提交记录
//	rs.GET("/submit-list", service.GetSubmitList)
//
//	// 获取测试案例
//	authAdmin.GET("/test-case", service.GetTestCase)
//
//	// 用户私有方法
//	authUser := rs.Group("/user", middlewares.AuthUserCheck())
//
//	return r
//}
