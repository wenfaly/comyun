package middleWares

import (
	"comyun/controller"
	"comyun/pkg/jwt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strings"
)

var TokenHeader = "Bearer"

const (
	KeyContextUserID = "user_id"
	KeyAccessControl = "access_role"
)

func JWTAuthMiddleWare() func(c *gin.Context){
	return func(c *gin.Context) {
		//获取请求头的认证信息并验证
		autnerHeader := c.Request.Header.Get("Authorization")
		if autnerHeader == ""{
			controller.ResponseError(c,controller.CodeNeedLogIn)
			c.Abort()
			return
		}
		//验证后获取token
		parts := strings.SplitN(autnerHeader," ",2)
		if !(len(parts) == 2 && parts[0] == TokenHeader) {
			controller.ResponseError(c,controller.CodeNeedLogIn)
			c.Abort()
			return
		}

		//对token进行解析
		//todo:后续如果需要进行token进行其他权限校验，进行接收并验证
		mc,err := jwt.ParseToken(parts[1])
		if err != nil{
			controller.ResponseError(c,controller.CodeNeedLogIn)
			zap.L().Error(err.Error())
			c.Abort()
			return
		}

		c.Set(KeyContextUserID,mc.UserID)
		c.Set(KeyAccessControl,mc.Access)
		c.Next()
	}
}
