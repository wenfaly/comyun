package models

import "time"

type PostParams struct {
	PostName string `json:"post_name"`//表单名称
	Description string `json:"description"`//描述
	CreateTime time.Time `json:"create_time"`//创建时间
	EndTime time.Time `json:"end_time"`//截止时间
}

type PostPublishParams struct {
	PostId int64 `json:"post_id"`
	Users []int64 `json:"users"`
}

type PostListSetParams struct {
	UserID int64 `json:"user_id"`
	Num int `json:"num"`
	Page int `json:"page"`
}

type PostListParams struct {
	PostID int64 `json:"post_id"`
	Num int `json:"num"`
	Page int `json:"page"`
}
