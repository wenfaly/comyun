package controller

import "github.com/gin-gonic/gin"

//用于编写错误码并将错误信息返回

type ResponseData struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func ResponseError(c *gin.Context, code ResCode) {
	rd := &ResponseData{
		Code: code,
		//依据状态码获取错误信息
		Msg:  code.GetMsg(),
		Data: nil,
	}
	c.JSON(200, rd)
}

func ResponseErrorWithMsg(c *gin.Context, code ResCode, data interface{}) {
	rd := &ResponseData{
		Code: code,
		//依据状态码获取错误信息
		Msg:  code.GetMsg(),
		Data: data,
	}
	c.JSON(200, rd)
}

func ResponseSuccess(c *gin.Context, data interface{}) {
	rd := &ResponseData{
		Code: CodeSuccess,
		//依据状态码获取错误信息
		Msg:  CodeSuccess.GetMsg(),
		Data: data,
	}
	c.JSON(200, rd)
}
