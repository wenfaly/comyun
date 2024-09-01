package routes

import (
	"comyun/controller"
	"comyun/logger"
	"comyun/middleWares"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetUp() *gin.Engine{
	r := gin.New()
	//todo: 限流中间件
	r.Use(logger.GinLogger(),logger.GinRecovery(true))

	//静态文件访问
	r.Static("/static","./static")

	v1 := r.Group("/api/v1")

	//注册用户
	//发送验证码
	v1.POST("/send_email",controller.SendEmail)
	//注册并填写信息
	v1.POST("/signup",controller.Signup)
	v1.POST("/judge_code",controller.JudgeCode)
	//登录并返回token
	v1.POST("/login_code",controller.LoginCode)
	v1.POST("/login_pass",controller.LoginPass)

	//认证JWT中间件
	v1.Use(middleWares.JWTAuthMiddleWare())
	{
		v2 := v1.Group("/department",middleWares.JWTAccessDepartmentMiddleWare())
		{
			//生成邀请链接
			v2.POST("/invite/link",controller.InviteLink)

			//批量导入,进行默认注册
			v2.GET("/invite/download",controller.InviteDownload)
			v2.POST("/invite/group",controller.InviteGroup)

			//公司登录日志
			v2.GET("/company/logger",controller.CompanyLogger)
			//成员统计
			v2.POST("/company/users",controller.CompanyUsers)
			//todo:权限设置

			//表单提交(保存，并不是发布)
			v2.POST("/post/set",controller.SetPost)
			//表单发布
			v2.POST("/post/publish",controller.PublishPost)
			//表单修改（权限）
			//修改表单组件和字段设置（修改源代码和组件设置）
			//todo:在mysql中修改表单的修改时间
			v2.POST("/post/adjust",controller.AdjustPost)
			//获取表单列表（权限）
			//设计的（发布的和未发布的）
			v2.POST("/post/list/set",controller.PostListSet)

			//统计表单数据
			v2.POST("/post/statistics",controller.StatisPost)
		}

		//新建公司
		v1.POST("/new_company",controller.NewCompany)

		////邀请链接
		//v1.POST("/invite/link",controller.InviteLink)
		////用户访问邀请链接，返回岗位信息
		////todo:此时未登录或注册先引导其登录注册
		//v1.GET("/invite/user",controller.InviteUser)
		////面对已登录用户，接收信息进行完善
		//v1.POST("/invite/login",controller.InviteLogin)
		////批量导入,进行默认注册
		//v1.GET("/invite/download",controller.InviteDownload)
		//v1.POST("/invite/group",controller.InviteGroup)

		//用户访问邀请链接，返回岗位信息
		//todo:此时未登录或注册先引导其登录注册
		v2.GET("/invite/user",controller.InviteUser)
		//面对已登录用户，接收信息进行完善
		v2.POST("/invite/login",controller.InviteLogin)

		//收到的
		v1.POST("/post/list/receive",controller.PostListReceive)
		//用户获取填写表单
		v1.POST("/user/post/detail",controller.GetUserPost)
		//提交表单数据
		v1.POST("/user/post/upload",controller.SetUserPost)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK,gin.H{
			"msg": "404",
		})
	})

	return r
}
