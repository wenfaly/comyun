package models

import "time"


type Post struct{
	ID int64 `db:"id"`//
	CompanyID int64 `db:"company_id"`//公司id
	PostName string `db:"post_name"`//表单名称
	Description string `db:"description"`//描述
	PostBy int64 `db:"post_by"`//发布者
	FieldId string `db:"field_id"`//字段组件设置id
	PublishIsNot bool `db:"publish_is_not"` //是否发布

	CreateTime time.Time `db:"create_time"`//创建时间
	UpdateTime time.Time `db:"update_time"`//更新时间
	EndTime time.Time `db:"end_time"`//截止时间
}

type PostTask struct {
	PostId int64 `db:"post_id"`//表单id
	PostTo int64 `db:"post_to"`//接收者
	Status bool `db:"status"`//完成情况
	CompleteTime time.Time `db:"complete_time"`//完成时间
}

type PostList struct {
	PostId int64 `json:"post_id" db:"id"`
	PublishIsNot bool `json:"publish_is_not" db:"publish_is_not"`
	PostName string `json:"post_name" db:"post_name"`//表单名称
	Description string `json:"description" db:"description"`//描述

	CreateTime time.Time `json:"create_time" db:"create_time"`//创建时间
	EndTime time.Time `json:"end_time" db:"end_time"`//截止时间
}
