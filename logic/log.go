package logic

import (
	"comyun/dao/mysql"
	"go.uber.org/zap"
)

func SetLoggerEmail(email string) error {
	log,err := mysql.GetUserFromEmail(email)
	if err != nil{
		zap.L().Error("get user err"+err.Error())
		return err
	}

	err = mysql.SetLog(log)

	return err
}

func SetLoggerTele(tele string) error {
	log,err := mysql.GetUserFromTele(tele)
	if err != nil{
		zap.L().Error("get user err"+err.Error())
		return err
	}

	err = mysql.SetLog(log)

	return err
}
