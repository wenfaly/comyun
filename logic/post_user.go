package logic

import (
	"comyun/dao/mongo"
	"comyun/dao/mysql"
	"comyun/models"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

func GetUserPost(postId int64)(fields *models.FormFieldGroup,err error) {
	//从post获取ObjectID
	objectID := ""
	objectID,err = mysql.GetObjectID(postId)
	if err != nil{
		zap.L().Error("get ObjectID error,"+err.Error())
		return nil,err
	}

	fields,err = mongo.GetFields(objectID)
	if err != nil{
		zap.L().Error("get PostFields error,"+err.Error())
		return nil,err
	}

	return fields,nil
}

// 对用户提交的表单进行参数校验
func UserFieldIfMatch(fg *models.UserFieldGroup)(err error,ok bool){
	//判断各个组件的值和类型是否匹配
	//fg.Fields刚好在接收时接收为[]struct对应传递的json字符串
	for _,v := range fg.Fields {
		//根据传递的值和类型进行断言，判断是否匹配成功
		ok = validateTypeAndValue(v)
		if !ok {
			return nil,ok
		}
	}

	//判断传入的值是否符合表单预定义的格式
	//根据post_id取出objectID对应的预定义
	oid := ""
	if oid,err = mysql.GetObjectID(fg.PostID);err != nil{
		zap.L().Error("get objectID err,"+err.Error())
		return err,false
	}
	fields := new(models.FormFieldGroup)
	if fields,err = mongo.GetFields(oid);err != nil{
		zap.L().Error("get postFields err,"+err.Error())
		return err,false
	}

	//判断值和预定义是否匹配
	if ok = validateDefineAndValue(fields.Fields,fg.Fields); !ok {
		zap.L().Error("validateDefineAndValue type and value error")
		return nil,false
	}

	return nil,true
}

//判断类型与值是否匹配
func validateTypeAndValue(field models.UserField) bool {
	switch field.Type {
	case "int":
		_,ok := field.Contents.(int)
		return ok
	case "int64":
		_,ok := field.Contents.(int64)
		return ok
	case "bool":
		_,ok := field.Contents.(int64)
		return ok
	case "FieldPosition":
		result := new(models.FieldPosition)
		err := mapstructure.Decode(field.Contents,result)
		if err != nil{
			zap.L().Error("mapstructure.Decode FieldPosition error,"+err.Error())
			return false
		}
	case "FieldAddress":
		result := new(models.FieldAddress)
		err := mapstructure.Decode(field.Contents,result)
		if err != nil{
			zap.L().Error("mapstructure.Decode FieldAddress error,"+err.Error())
			return false
		}
	case "FieldCard":
		result := new(models.FieldCard)
		err := mapstructure.Decode(field.Contents,result)
		if err != nil{
			zap.L().Error("mapstructure.Decode FieldCard error,"+err.Error())
			return false
		}
	case "string":
		_,ok := field.Contents.(string)
		return ok
	default:
		return false
	}

	return true
}

func validateDefineAndValue(formFields []models.FormField,fg []models.UserField) bool {
	for i,v := range formFields {
		if v.Name != fg[i].Name {
			zap.L().Error("validateDefineAndValue error :Name not match")
			return false
		}
		//仅对特殊类型进行校验
		switch v.Type {
		case "Address":
			if !validateAddressDefined(v,fg[i]) {
				zap.L().Error("error to assertion Address")
				return false
			}
			continue
		case "PersonCard":
			if !validatePersonCardDefined(v,fg[i]) {
				zap.L().Error("error to assertion PersonCard")
				return false
			}
			continue
		case "Position":
			if !validatePositionDefined(v,fg[i]) {
				zap.L().Error("error to assertion Position")
				return false
			}
			continue
		case "int64": continue
		case "int": continue
		case "bool":continue
		case "string":continue
		default:
			return false
		}
	}

	return true
}

func validateAddressDefined(ff models.FormField,uf models.UserField) bool {
	addressD,ok := ff.OwnFunc.(primitive.D)
	if !ok {
		zap.L().Error("validateDefineAndValue error : assertion Address primitive.D{} failed")
		return false
	}

	addressUser := new(models.FieldAddress)
	//mapstructure.Decode 的反射匹配有时不能完全解析，需要在被接收的对象加上 `mapstructure` 标签
	err := mapstructure.Decode(uf.Contents,addressUser)
	if err != nil {
		zap.L().Error("validateDefineAndValue error : FieldAddress assertion failed,"+err.Error())
		return false
	}

	//判断是否需要详细地址
	v,ok := addressD[0].Value.(bool)
	if v {
		if addressUser.AddressDetail != ""{
			return true
		}
		return false
	}else{
		if addressUser.AddressDetail == ""{
			return true
		}
		return false
	}
}

func validatePersonCardDefined(ff models.FormField,uf models.UserField) bool {
	pCard := new(models.PersonCard)
	err := mapstructure.Decode(ff.OwnFunc,pCard)
	if err != nil {
		zap.L().Error("validateDefineAndValue error :PersonCard assertion failed,"+err.Error())
		return false
	}

	card := new(models.FieldCard)
	err = mapstructure.Decode(uf.Contents,card)
	if err != nil {
		zap.L().Error("validateDefineAndValue error :FieldCard assertion failed")
		return false
	}

	//检验证件上的姓名
	if pCard.CardName {
		if card.CardName == "" {
			return false
		}
	}else{
		if card.CardName != "" {
			return false
		}
	}
	//检验证件上的性别
	if pCard.CardGender {
		if card.CardGender == "" {
			return false
		}
	}else{
		if card.CardGender != "" {
			return false
		}
	}
	//检验证件上的证件号
	if pCard.CardID {
		if card.CardID == "" {
			return false
		}
	}else{
		if card.CardID != "" {
			return false
		}
	}
	//检验证件上的生日
	if pCard.CardBirth {
		if card.CardBirth == "" {
			return false
		}
	}else{
		if card.CardBirth != "" {
			return false
		}
	}

	return true
}

func validatePositionDefined(ff models.FormField,uf models.UserField) bool {
	pPosition,ok := ff.OwnFunc.(primitive.D)
	if !ok {
		zap.L().Error("validateDefineAndValue error : assertion FormField primitiveD assertion failed")
		return false
	}

	position := new(models.FieldPosition)
	err := mapstructure.Decode(uf.Contents,position)
	if err != nil {
		zap.L().Error("validateDefineAndValue error :FieldPosition assertion failed")
		return false
	}

	//检验位置信息
	v,ok := pPosition[0].Value.(bool)
	if v {
		if !(position.CenterLimit && position.CenterAcross) {
			return false
		}
	}
	return true
}

