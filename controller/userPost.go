package controller

import (
	"comyun/dao/mongo"
	"comyun/logic"
	"comyun/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

//用户获取需要填写的表单
func GetUserPost(c *gin.Context) {
	pidStr,ok := c.GetQuery(KeyPostID)
	if !ok {
		zap.L().Error("error to get param post_id")
		ResponseError(c,CodeInValidParam)
		return
	}
	pid,err := strconv.ParseInt(pidStr,10,64)
	if err != nil{
		zap.L().Error("error to parseInt post_id,"+err.Error())
		ResponseError(c,CodeInValidParam)
		return
	}

	postFields,err := logic.GetUserPost(pid)
	if err != nil{
		zap.L().Error("error to get UserPost,"+err.Error())
		ResponseError(c,CodeServerBusy)
		return
	}

	ResponseSuccess(c,postFields)
}

func SetUserPost(c *gin.Context) {
	uf := new(models.UserFieldGroup)
	if err := c.ShouldBindJSON(uf);err != nil{
		zap.L().Error("ShouldBindJson models.UserFieldGroup error,"+err.Error())
		ResponseError(c,CodeInValidParam)
		return
	}

	//参数校验
	if err,ok := logic.UserFieldIfMatch(uf);err != nil || !ok {
		if err != nil{
			zap.L().Error("UserField not match,"+err.Error())
			ResponseError(c,CodeServerBusy)
			return
		}
		zap.L().Error("UserField not match")
		ResponseError(c,CodeInValidParam)
		return
	}

	//数据进行保存
	if err := mongo.SaveUserFields(uf);err != nil{
		zap.L().Error("save UserField error,"+err.Error())
		ResponseError(c,CodeServerBusy)
		return
	}

	ResponseSuccess(c,nil)
}

func StatisPost(c *gin.Context) {
	pl := new(models.PostListParams)
	if err := c.ShouldBindJSON(pl);err != nil{
		zap.L().Error("Get post_id error,"+err.Error())
		ResponseError(c,CodeInValidParam)
		return
	}

	posts,err := mongo.GetUserPosts(pl)
	if err != nil{
		zap.L().Error("mongo.GetUserPosts error,"+err.Error())
		ResponseError(c,CodeServerBusy)
		return
	}

	ResponseSuccess(c,posts)
}
