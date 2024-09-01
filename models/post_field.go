package models

// mongoDB中存储结构
type FormFieldGroup struct {
	CompanyID int64 `bson:"company_id"`
	UserID int64 `bson:"user_id"`//发布用户id
	Role int `bson:"role"`//1 为发起者，0 为提交者
	Code string `bson:"code"`
	Fields []FormField `bson:"fields"`
}

//后端接收示例
//type Fields []FormField

//表单组件的接收
type FormField struct{
	Name string `bson:"name" json:"name"`//标题（用户设置的）
	Detail string `bson:"detail" json:"detail"`//说明
	Type string `bson:"type" json:"type"`//内容的类型（与后端相同的类型，包含特殊设置类型）
	Required bool `bson:"require" json:"require"`//必填

	OwnFunc interface{} `bson:"own_func" json:"own_func"`//特殊设置
}

// 以下为组件的特殊功能

type Address struct {
	AddressType bool `bson:"address_type" json:"address_type"`//填写是否包含详细地址
}

type PersonCard struct {
	CardName bool `bson:"card_name" json:"card_name"`//是否提取姓名
	CardBirth bool `bson:"card_birth" json:"card_birth"`//是否提取出生日期
	CardGender bool `bson:"card_gender" json:"card_gender"`//是否提取性别
	CardID bool `bson:"card_id" json:"card_id"`//是否提取身份证号码
}

type Position struct {
	Limit bool `bson:"limit" json:"limit"`//是否设置中心点附近范围内定位
	CenterLongitude string `bson:"center_longitude" json:"center_longitude"`//设置中心点——中心点经度
	CenterLatitude string `bson:"center_latitude" json:"center_latitude"`//设置中心点——中心点纬度
	Length int `bson:"length" json:"length"`//设置中心点——中心点距离限制
	Description string `bson:"description" json:"description"`//设置中心点——中心点位置详细地址
}
