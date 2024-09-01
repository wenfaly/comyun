package controller

type ResCode int64

const (
	CodeSuccess = ResCode(1000) + iota
	CodeInValidParam
	CodeInValidFile

	CodeUserExist
	CodeUserNotExist
	CodeInvalidPassword
	CodeServerBusy

	CodeAccessDenied
	CodeNeedLogIn
	CodeEmailNotMatch
	CodeEmailExist
	CodeEmailNotExist
	CodeEmailCodeDefault
	CodeInviteLinkOvertime
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:         "success",
	CodeInValidParam:    "请求参数错误",
	CodeInValidFile: "接收文件错误",
	CodeUserExist:       "用户已存在",
	CodeUserNotExist:    "用户不存在",
	CodeInvalidPassword: "用户名或密码错误",
	CodeAccessDenied: "用户权限不足",
	CodeNeedLogIn:       "需要登录",
	CodeEmailNotMatch: "邮箱格式错误",
	CodeEmailExist: "邮箱已存在",
	CodeEmailNotExist: "邮箱不存在",
	CodeEmailCodeDefault: "验证码错误",
	CodeInviteLinkOvertime: "邀请链接过期",
}

// GetMsg 通过错误码获取对应的错误信息
func (r ResCode) GetMsg() string {
	msg, ok := codeMsgMap[r]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}

	return msg
}
