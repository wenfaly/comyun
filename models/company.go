package models

type Company struct {
	Id int64 `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	CompanyId int64 `json:"company_id" db:"company_id"`
	OwnerId int64 `json:"owner_id" db:"owner_id"`
}

type ComUserParam struct{
	CompanyID int64 `json:"company_id"`
	Page int `json:"page"`
	Sum int `json:"sum"`
}

