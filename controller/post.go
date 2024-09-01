package controller

import (
	"comyun/dao/mongo"
	"comyun/dao/mysql"
	"comyun/logic"
	"comyun/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SetPost(c *gin.Context){
	strFields := c.PostForm(KeyFields)
	strCode := c.PostForm(KeyCode)
	strPost := c.PostForm(KeyPost)
	if strFields == "" || strCode == "" || strPost == ""{
		zap.L().Error("get key fields error")
		ResponseError(c,CodeInValidParam)
		return
	}

	//处理字段设置参数
	fields := new([]models.FormField)
	//todo:参数为空的情况
	if err := json.Unmarshal([]byte(strFields),fields);err != nil{
		zap.L().Error("fields unmarshal error")
		ResponseError(c,CodeInValidParam)
		return
	}

	//获取信息存入mongo
	i := c.GetInt64(KeyContextUserID)
	if i == 0{
		zap.L().Error("c.GetInt64 error")
		ResponseError(c,CodeServerBusy)
		return
	}
	id := mysql.GetCompanyID(i)
	if id == 0{
		zap.L().Error("GetCompanyID error")
		ResponseError(c,CodeServerBusy)
		return
	}
	fg := &models.FormFieldGroup{
		CompanyID: id,
		UserID:    i,
		Role:      1,
		Code: 	strCode,
		Fields:    *fields,
	}
	p := new(models.PostParams)
	if err := json.Unmarshal([]byte(strPost),p);err != nil{
		zap.L().Error("Unmarshal PostParams error"+err.Error())
		ResponseError(c,CodeInValidParam)
		return
	}
	post := &models.Post{
		CompanyID:   id,
		PostName:    p.PostName,
		Description: p.Description,
		PostBy:      i,
		PublishIsNot: false,
		CreateTime:  p.CreateTime,
		EndTime:     p.EndTime,
	}

	if err := mysql.SetPost(fg,post);err != nil{
		zap.L().Error("insert fieldGroup error"+err.Error())
		ResponseError(c,CodeServerBusy)
		return
	}

	ResponseSuccess(c,nil)
}

func AdjustPost(c *gin.Context) {
	strFields := c.PostForm(KeyFields)
	strCode := c.PostForm(KeyCode)
	strPostId := c.PostForm(KeyPostID)
	if strFields == "" || strCode == "" || strPostId == ""{
		zap.L().Error("get key fields error")
		ResponseError(c,CodeInValidParam)
		return
	}

	//处理参数
	fields := new([]models.FormField)
	//todo:参数为空的情况
	if err := json.Unmarshal([]byte(strFields),fields);err != nil{
		zap.L().Error("fields unmarshal error")
		ResponseError(c,CodeInValidParam)
		return
	}

	if err := mongo.AdjustFields(strFields,strPostId);err != nil{
		zap.L().Error("insert fieldGroup error"+err.Error())
		ResponseError(c,CodeServerBusy)
		return
	}

	ResponseSuccess(c,nil)
}

func PublishPost(c *gin.Context) {
	p := new(models.PostPublishParams)
	if err := c.ShouldBindJSON(p);err != nil{
		zap.L().Error("shouldBindJson PublishParams error"+err.Error())
		ResponseError(c,CodeInValidParam)
		return
	}

	if err := logic.PublishPost(p);err != nil{
		zap.L().Error("logic.PublishPost error"+err.Error())
		ResponseError(c,CodeServerBusy)
		return
	}

	ResponseSuccess(c,nil)
}

//获取所有设计的表单
func PostListSet(c *gin.Context) {
	p := new(models.PostListSetParams)
	if err := c.ShouldBindJSON(p);err != nil {
		zap.L().Error("get key KeyContextUserID error")
		ResponseError(c,CodeInValidParam)
		return
	}

	//获取设计的表单
	posts,err := mysql.GetPostSetList(p)
	if err != nil{
		zap.L().Error("error to get postSetList,"+err.Error())
		ResponseError(c,CodeServerBusy)
		return
	}

	ResponseSuccess(c,posts)
}

//获取所有收到的表单
func PostListReceive(c *gin.Context) {
	p := new(models.PostListSetParams)
	if err := c.ShouldBindJSON(p);err != nil {
		zap.L().Error("get key KeyContextUserID error")
		ResponseError(c,CodeInValidParam)
		return
	}

	posts,err := logic.GetPostReceiveList(p)
	if err != nil{
		zap.L().Error("error to get postReceiveList,"+err.Error())
		ResponseError(c,CodeServerBusy)
		return
	}

	ResponseSuccess(c,posts)
}
