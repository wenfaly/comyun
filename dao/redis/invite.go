package redis

import (
	"comyun/models"
	"encoding/json"
	"time"
)

func InviteLink(invite *models.UserInvite,code string) error {
	i,err := json.Marshal(invite)
	if err != nil{
		return err
	}
	client.Set(
		getRedisKey(KeyInviteCode+code),i,
		time.Duration(invite.Duration)*time.Hour)

	return nil
}

func GetInviteCode(code string)([]byte,error){
	v,err := client.Get(getRedisKey(KeyInviteCode+code)).Result()
	if err != nil{
		if err == Nil{
			return nil,nil
		}
		return nil,err
	}

	return []byte(v),nil
}
