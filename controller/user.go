package controller

import (
	"comyun/dao/mysql"
	"comyun/dao/redis"
	"comyun/logic"
	"comyun/models"
	"comyun/pkg/jwt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Signup(c *gin.Context){
	//接收用户的完善信息
	u := new(models.User)
	if err := c.ShouldBindJSON(u);err != nil{
		zap.L().Error("shouldBindJson error"+err.Error())
		ResponseError(c,CodeInValidParam)
		return
	}

	//参数的校验

	//注册人员
	err := logic.SignupUser(u)
	if err != nil{
		zap.L().Error("signupUser err"+err.Error())
		ResponseError(c,CodeServerBusy)
		return
	}

	//返回响应
	ResponseSuccess(c,nil)
}

// LoginCode 验证码登录
func LoginCode(c *gin.Context){
	cp := new(models.CodeParams)
	if err := c.ShouldBindJSON(cp);err != nil{
		zap.L().Info("ShouldBindJson error"+err.Error())
		ResponseError(c,CodeInValidParam)
		return
	}

	ok,err := redis.JudgeCode(cp)
	if err != nil{
		zap.L().Error("JudgeCode error"+err.Error())
		ResponseError(c,CodeServerBusy)
		return
	}
	if !ok{
		ResponseError(c,CodeEmailCodeDefault)
		return
	}

	//获取user_id生成token
	id,err := mysql.GetUserID(cp.Email)
	if err != nil{
		zap.L().Error("get id err"+err.Error())
		ResponseError(c,CodeEmailExist)
		return
	}
	token,err := jwt.NewToken(id)
	if err != nil{
		zap.L().Error("new token err"+err.Error())
		ResponseError(c,CodeServerBusy)
		return
	}

	//todo:登陆成功，记录日志
	logic.SetLoggerEmail(cp.Email)

	ResponseSuccess(c,token)
}

// LoginPass 密码登录
func LoginPass(c *gin.Context){
	up := new(models.UserParams)
	if err := c.ShouldBindJSON(up);err != nil{
		zap.L().Error("shouldBindJson err"+err.Error())
		ResponseError(c,CodeInValidParam)
		return
	}

	id,err := logic.LoginPass(up)
	if err != nil{
		zap.L().Error("loginPass err"+err.Error())
		if err == mysql.ErrorPassword{
			ResponseError(c,CodeInvalidPassword)
			return
		}
		ResponseError(c,CodeServerBusy)
		return
	}

	token,err := jwt.NewToken(id)
	if err != nil{
		zap.L().Error("new token err"+err.Error())
		ResponseError(c,CodeServerBusy)
		return
	}

	//todo:登陆成功，记录日志
	logic.SetLoggerTele(up.Telephone)

	ResponseSuccess(c,token)
}
