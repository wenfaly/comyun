package mysql

import (
	"comyun/dao/mongo"
	"comyun/models"
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"strings"
)

func SetPost(fields *models.FormFieldGroup,post *models.Post) error {
	//存入数据
	//mysql事务
	tx,err := db.Begin()
	if err != nil {
		zap.L().Error("db.Begin error"+err.Error())
		return err
	}

	sqlstr := `insert into post(
		company_id,post_name,description,post_by,publish_is_not,create_time,end_time) 
		values (?,?,?,?,?,?)`
	result,err := tx.Exec(sqlstr,post.CompanyID,post.PostName,post.Description,post.PostBy,post.PublishIsNot,post.CreateTime,post.EndTime)
	if err != nil{
		zap.L().Error("insert post error"+err.Error())
		tx.Rollback()
		return err
	}

	//获取field_id
	r,err := mongo.CollPost.InsertOne(context.TODO(),*fields)
	if err != nil{
		zap.L().Error("Insert post error"+err.Error())
		tx.Rollback()
		return err
	}

	//获取插入文档的ObjectID并转换格式
	fid := r.InsertedID.(primitive.ObjectID).String()
	fid = fid[10:len(fid)-2]

	//更新post中的ObjectID(field_id)
	sqlstr = "update post set field_id = ? where id = ?"
	id,err := result.LastInsertId()
	_,err = tx.Exec(sqlstr,fid,id)
	if err != nil{
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil{
		return err
	}
	return nil
}

func InsertPostTasks(params []interface{},queryStr []string,postID int64) error {
	tx,err := db.Begin()
	if err != nil{
		zap.L().Error("db.begin error"+err.Error())
		return err
	}
	//插入数据PostTask，更改Post状态
	sqlstr := fmt.Sprintf(`insert into post_task (
		post_id,post_to) values %s`,strings.Join(queryStr,","))
	_,err = tx.Exec(sqlstr,params...)
	if err != nil{
		zap.L().Error("insert post_task error,"+err.Error())
		tx.Rollback()
		return err
	}

	sqlstr = `update post set publish_is_not = ? where id = ?`
	_,err = tx.Exec(sqlstr,1,postID)
	if err != nil{
		zap.L().Error("update post publish_exist error"+err.Error())
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func GetPostSetList(params *models.PostListSetParams)(posts []models.PostList,err error){
	sqlstr := `select 
    	id,post_name,description,publish_is_not,create_time,end_time 
		from post where post_by = ? limit ? offset ?`

	err = db.Select(&posts,sqlstr,params.UserID,params.Num,(params.Page-1)*params.Num)
	if err != nil{
		zap.L().Error("error to select post")
		return nil,err
	}

	return
}

func GetPostIds(params *models.PostListSetParams)(ids []int64,err error){
	sqlstr := `select post_id from post_task where post_to = ? limit ? offset ?`

	err = db.Select(&ids,sqlstr,params.UserID,params.Num,(params.Page-1)*params.Num)
	if err != nil{
		zap.L().Error("select post_id error,"+err.Error())
		return nil,err
	}

	return
}

func GetPostReceiveList(ids []int64)(posts []models.PostList,err error){
	sqlstr := `select id,post_name,description,publish_is_not,create_time,end_time 
		from post where id IN(?);`

	query,args,err := sqlx.In(sqlstr,ids)
	if err != nil{
		zap.L().Error("select sqlx.In error,"+err.Error())
		return nil,err
	}

	query = db.Rebind(query)
	err = db.Select(&posts,query,args...)
	if err != nil{
		zap.L().Error("db.Select error,"+err.Error())
		return nil,err
	}

	return
}

func GetObjectID(pid int64)(objectID string,err error){
	sqlstr := `select field_id from post where id = ?`

	err = db.Get(&objectID,sqlstr,pid)
	if err != nil{
		zap.L().Error("get ObjectID error,"+err.Error())
		return "",err
	}

	return objectID,nil
}
