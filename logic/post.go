package logic

import (
	"comyun/dao/mysql"
	"comyun/models"
	"go.uber.org/zap"
)

const queryString = "(?,?)"

//表单对成员进行发布
func PublishPost(p *models.PostPublishParams) error {
	n := len(p.Users)
	params := make([]interface{},0,n*2)
	queryStr := make([]string,0,n*2)
	for _,v := range p.Users{
		queryStr = append(queryStr,queryString)
		params = append(params,p.PostId)
		params = append(params,v)
	}

	return mysql.InsertPostTasks(params,queryStr,p.PostId)
}

func GetPostReceiveList(p *models.PostListSetParams)(posts []models.PostList,err error){
	//获取个人收到的post_id
	var ids []int64
	ids,err = mysql.GetPostIds(p)
	if err != nil{
		zap.L().Error("mysql.GetPostIds error,"+err.Error())
		return nil,err
	}

	posts,err = mysql.GetPostReceiveList(ids)
	if err != nil{
		zap.L().Error("mysql.GetPostReceiveList error,"+err.Error())
		return nil,err
	}

	return
}
