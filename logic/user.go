package logic

import (
	"comyun/dao/mysql"
	"comyun/models"
	"comyun/pkg/snowflake"
)

//注册时进行生成id，存入数据库
func SignupUser(u *models.User) error {
	u.UserID = snowflake.GenID()
	return mysql.SignupUser(u)
}

func LoginPass(up *models.UserParams)(int64,error){
	return mysql.LoginPass(up)
}
