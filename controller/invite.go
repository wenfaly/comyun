package controller

import (
	"comyun/dao/mysql"
	"comyun/dao/redis"
	"comyun/logic"
	"comyun/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/exp/rand"
	"os"
	"strconv"
	"time"
)

// 接收包含设置岗位等信息
func InviteLink(c *gin.Context){
	invite := new(models.UserInvite)
	if err := c.ShouldBindJSON(invite);err != nil{
		zap.L().Error("shouldBindJson userInvite err"+err.Error())
		ResponseError(c,CodeInValidParam)
		return
	}

	//参数校验

	//生成随机校验参数
	rand.Seed(uint64(time.Now().UnixNano()))
	code := strconv.Itoa(10000000+rand.Intn(89999999))
	//在redis存储验证参数和设置信息的对应关系
	if err := redis.InviteLink(invite,code);err != nil{
		zap.L().Error(" userInvite err"+err.Error())
		ResponseError(c,CodeServerBusy)
		return
	}

	ResponseSuccess(c,InviteLinkPrefix+code)
}

func InviteUser(c *gin.Context){
	code,ok := c.GetQuery("invite")
	if !ok{
		zap.L().Error("query \"invite\" error")
		ResponseError(c,CodeInValidParam)
		return
	}

	invite,err := redis.GetInviteCode(code)
	if err != nil{
		zap.L().Error("get redis key"+err.Error())
		ResponseError(c,CodeServerBusy)
		return
	}

	if invite == nil {
		ResponseError(c,CodeInviteLinkOvertime)
		return
	}

	i := new(models.UserInvite)
	if err = json.Unmarshal(invite,i);err != nil{
		zap.L().Error("Unmarshal inviteCode err"+err.Error())
		ResponseError(c,CodeServerBusy)
		return
	}

	ResponseSuccess(c,i)
}

func InviteLogin(c *gin.Context){
	i := new(models.UserInvite)
	err := c.ShouldBindJSON(i)
	if err != nil{
		zap.L().Error("shouldBindJson userInvite err"+err.Error())
		ResponseError(c,CodeInValidParam)
		return
	}

	//todo:检测是否加入其他企业

	if err = mysql.InviteLogin(i);err != nil{
		zap.L().Error("mysql InviteLogin err"+err.Error())
		ResponseError(c,CodeServerBusy)
		return
	}

	ResponseSuccess(c,nil)
}

func InviteDownload(c *gin.Context){
	filepath := "./uploads/xlsx/example.xlsx"
	if _,err := os.Stat(filepath);os.IsNotExist(err){
		zap.L().Error("example xlsx file not find"+err.Error())
		ResponseError(c,CodeServerBusy)
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename=example.txt")
	c.Header("Content-Type", "application/octet-stream")
	c.File(filepath)
}

func InviteGroup(c *gin.Context){
	cid := c.PostForm("cid")
	if cid == ""{
		zap.L().Error("get cid err")
		ResponseError(c,CodeInValidParam)
		return
	}
	handler,err := c.FormFile("file")
	if err != nil{
		zap.L().Error("get file err"+err.Error())
		ResponseError(c,CodeInValidFile)
		return
	}

	//打开文件流，不进行保存后再打开
	file,err := handler.Open()
	if err != nil{
		zap.L().Error("cannot open file"+err.Error())
		ResponseError(c,CodeInValidFile)
		return
	}
	defer file.Close()

	//处理xlsx文件并进行保存
	if err = logic.InviteGroup(cid,file);err != nil{
		zap.L().Error("inviteGroup error"+err.Error())
		ResponseError(c,CodeServerBusy)
		return
	}

	ResponseSuccess(c,nil)
}

