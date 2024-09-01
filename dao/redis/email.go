package redis

import (
	"comyun/models"
	"fmt"
	"go.uber.org/zap"
	"time"
)

const emailOvertime = 10*time.Minute

func SetEmail(email string,code int)error{
	pipeline := client.TxPipeline()

	pipeline.Set(getRedisKey(KeyEmailSet+email),code, emailOvertime)

	_,err := pipeline.Exec()
	return err
}

func JudgeCode(cp *models.CodeParams) (bool,error) {
	cmd := client.Get(getRedisKey(KeyEmailSet+cp.Email))
	code,err := cmd.Int()
	fmt.Println(err)
	if err != nil{
		if err != Nil{
			return false,nil
		}
		zap.L().Error("get code err"+err.Error())
		return false,err
	}

	if code != cp.Code{
		return false,nil
	}
	return true,nil
}
