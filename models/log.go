package models

import "time"

type Log struct {
	Id int64 `db:"id"`
	UserId int64 `db:"user_id"`
	UserName string `db:"user_name"`
	CompanyID string `db:"company_id"`
	Department string `db:"department"`
	Role string `db:"role"`
	LogTime time.Time `db:"log_time"`
}

type ComLogParam struct{
	UserId int64 `json:"user_id"`
	UserName string `json:"user_name"`
	StartTime time.Time `json:"start_time"`
	EndTime time.Time `json:"end_time"`
	CompanyID string `json:"company_id"`
	Page int `json:"page"`
	Sum int `json:"sum"`
}

//返回给用户的消息
type ComLogMess struct {
	UserId int64 `json:"user_id"`
	UserName string `json:"user_name"`
	Department string `json:"department"`
	Role string `json:"role"`
	LogTime time.Time `json:"log_time"`
}
