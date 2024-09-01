package controller

import (
	"comyun/dao/mysql"
	"comyun/logic"
	"comyun/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// NewCompany默认本人为创建者
func NewCompany(c *gin.Context){
	//接收公司的完善信息
	p := new(models.Company)
	if err := c.ShouldBindJSON(p);err != nil{
		zap.L().Error("company shouldBindJson error"+err.Error())
		ResponseError(c,CodeInValidParam)
		return
	}

	//注册公司
	err := logic.NewCompany(p)
	if err != nil{
		zap.L().Error("signupUser err"+err.Error())
		ResponseError(c,CodeServerBusy)
		return
	}

	//返回响应
	ResponseSuccess(c,nil)
}

//获取公司在职人员
func CompanyUsers(c *gin.Context){
	cu := new(models.ComUserParam)
	if err := c.ShouldBindJSON(cu);err != nil{
		zap.L().Error("shouldBindJson err"+err.Error())
		ResponseError(c,CodeInValidParam)
		return
	}

	users,err := mysql.GetUsers(cu)
	if err != nil{
		zap.L().Error("get users err"+err.Error())
		ResponseError(c,CodeServerBusy)
		return
	}

	ResponseSuccess(c,users)
}

//根据不同的查询条件获取登录日志
func CompanyLogger(c *gin.Context){
	cl := new(models.ComLogParam)
	err := c.ShouldBindJSON(cl)
	if err != nil{
		zap.L().Error("CompanyLogger error"+err.Error())
		ResponseError(c,CodeInValidParam)
		return
	}

	//对参数进行校验选择不同的查询函数
	//登录时间查询
	var logs []models.ComLogMess
	if !cl.StartTime.IsZero() && !cl.EndTime.IsZero(){
		if cl.UserId == 0{
			logs,err = mysql.GetLoggerTime(cl)
		}else{
			logs,err = mysql.GetLoggerPT(cl)
		}
	}else{
		if cl.UserId == 0{
			logs,err = mysql.GetLoggerPerson(cl)
		}else{
			logs,err = mysql.GetLoggerPage(cl)
		}
	}

	if err != nil{
		zap.L().Error("mysql.GetLogger error"+err.Error())
		ResponseError(c,CodeServerBusy)
		return
	}

	ResponseSuccess(c,logs)
}
