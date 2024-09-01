package models

import "time"

type UserFieldGroup struct {
	CompanyID int64 `bson:"company_id" json:"company_id"`
	PostID int64 `bson:"post_id" json:"post_id"`//表单设置的post_id
	UserID int64 `bson:"user_id" json:"user_id"`//提交者id
	Role bool `bson:"role" json:"role"`
	SubmissionDate time.Time `bson:"submission_date" json:"submission_date"`//提交日期

	Fields []UserField `bson:"fields" json:"fields"`
}

type UserField struct {
	Name string `bson:"name" json:"name"`//标题
	Type string `bson:"type" json:"type"`//内容的类型（与后端相同的类型，包含特殊设置类型）

	Contents interface{} `bson:"contents" json:"contents"`//内容,interface{}兼容对象和常用类型
}

// 定位对象
type FieldPosition struct {
	Area string `bson:"area" json:"area" mapstructure:"area"`//省市区
	Description string `bson:"description" json:"description" mapstructure:"description"`//详细地址描述
	Longitude string `bson:"longitude" json:"longitude" mapstructure:"longitude"`//定位位置经度
	Latitude string `bson:"latitude" json:"latitude" mapstructure:"latitude"`//定位位置纬度

	CenterLimit bool `bson:"center_limit" json:"center_limit" mapstructure:"center_limit"`//是否设置中心点位置判断
	//todo:是否需要将中心点的位置判断单独写一个接口，在后端进行判断
	CenterAcross bool `bson:"center_across" json:"center_across" mapstructure:"center_across"`//当设置中心点时，是否满足范围
}

//地址
type FieldAddress struct {
	Area string `bson:"area" json:"area" mapstructure:"area"`//省市区
	AddressDetail string `bson:"address_detail" json:"address_detail" mapstructure:"address_detail"`//详细地址
}

// 文字识别后的身份证
type FieldCard struct {
	CardURL string `bson:"card_url" json:"card_url" mapstructure:"card_url"`//身份证照片的存储地址
	CardName string `bson:"card_name" json:"card_name" mapstructure:"card_name"`//姓名
	CardBirth string `bson:"card_birth" json:"card_birth" mapstructure:"card_birth"`//出生日期
	CardGender string `bson:"card_gender" json:"card_gender" mapstructure:"card_gender"`//性别
	CardID string `bson:"card_id" json:"card_id" mapstructure:"card_id"`//身份证号码
}
