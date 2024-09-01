package controller

import (
	"comyun/dao/redis"
	"comyun/logic"
	"comyun/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//SendEmail 将验证码返回给前端
func SendEmail(c *gin.Context) {
	//todo: 防止恶意发送邮件，每60秒允许发送一次(前端)
	//获取邮件信息
	ep := new(models.EmailParams)
	if err := c.ShouldBindJSON(ep); err != nil {
		zap.L().Error("error to get email")
		ResponseError(c, CodeInValidParam)
		return
	}

	//todo: 验证邮箱是否存在(注册-不存在，登录-存在)
	//验证登录模式
	if ep.Model != "signup" && ep.Model != "login" {
		ResponseError(c, CodeInValidParam)
		return
	}
	//邮箱格式是否符合规范
	if err := logic.EmailMatch(ep.Email); err != nil {
		if err == logic.ErrorNotMatch {
			zap.L().Error("error to match email")
			ResponseError(c, CodeEmailNotMatch)
			return
		}
		zap.L().Error("error to match email")
		ResponseError(c, CodeServerBusy)
		return
	}
	//验证邮箱是否存在
	if exist := logic.EmailExist(ep); !exist {
		zap.L().Info("user email model error")
		if !exist {
			if ep.Model == "signup"{
				ResponseError(c,CodeEmailExist)
			}else{
				ResponseError(c,CodeEmailNotExist)
			}
			return
		}
	}

	code, err := logic.SendEmail(ep.Email)
	//todo: 原生和三方库发送逻辑有什么不同，原生方法为什么错误
	//code,err := logic.SendEmail(email)
	if err != nil {
		zap.L().Error("error to seng email" + err.Error())
		ResponseError(c, CodeServerBusy)
		return
	}

	//todo: 将验证码及时间存入redis
	if err = redis.SetEmail(ep.Email, code); err != nil {
		zap.L().Error("error to set email & code" + err.Error())
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, code)
}

//注册时检验验证码
func JudgeCode(c *gin.Context){
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

	ResponseSuccess(c,nil)
}


