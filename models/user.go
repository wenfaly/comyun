package models

//注册、个人信息表
type User struct {
	Id int64 `json:"id" db:"id"`
	UserID int64 `json:"user_id" db:"user_id"`
	Boss bool `json:"boss" db:"boss"`//boss-1,employee-0
	CompanyID string `json:"company_id" db:"company_id"`
	Name string `json:"name" db:"name"`
	Gender bool `json:"gender" db:"gender"`
	Telephone string `json:"telephone" db:"telephone"`
	Email string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
	Department string `json:"department" db:"department"`
	Role string `json:"role" db:"role"`
}

//登录时进行的参数校验
type UserParams struct{
	Telephone string `json:"telephone"`
	Password string `json:"password"`
}

type UserInvite struct{
	UserID int64 `json:"user_id"`
	Name string `json:"name" binding:"required"`
	CompanyID int64 `json:"company_id" binding:"required"`
	Department string `json:"department"`
	Role string `json:"role"`
	Duration int `json:"duration"`//单位：小时
}
