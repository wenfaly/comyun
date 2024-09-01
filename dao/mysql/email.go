package mysql

import "go.uber.org/zap"

func EmailExist(email string)(exist bool){
	sqlStr := "select count(user_id) from user where email = ?"

	count := 0
	if err := db.Get(&count,sqlStr,email);err != nil{
		zap.L().Error(err.Error())
		return false
	}

	if count > 0 {
		return true
	}

	return false
}
