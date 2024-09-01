package middleWares

import (
	"comyun/controller"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const RoleStrange = 10

// 权限校验-对本部门表单的查看、编辑，表单的发布
func JWTAccessDepartmentMiddleWare() func(c *gin.Context) {
	return func(c *gin.Context) {
		v,exist := c.Get(KeyAccessControl)
		role := v.(int)
		if !exist || role < 2 {
			zap.L().Error("error to access next middleWare")
			controller.ResponseError(c,controller.CodeAccessDenied)
			c.Abort()
			return
		}

		c.Next()
	}
}

// 权限校验-对本部门表单的查看、编辑，表单的发布
func JWTAccessCEOMiddleWare() func(c *gin.Context) {
	return func(c *gin.Context) {
		v,exist := c.Get(KeyAccessControl)
		role := v.(int)
		if !exist || role < 3 {
			zap.L().Error("error to access next middleWare")
			controller.ResponseError(c,controller.CodeAccessDenied)
			c.Abort()
			return
		}

		c.Next()
	}
}
