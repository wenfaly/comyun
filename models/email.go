package models

type EmailParams struct {
	Email string`json:"email"`
	Model string `json:"model"`
}

type CodeParams struct {
	Email string `json:"email"`
	Code int `json:"code"`
}
