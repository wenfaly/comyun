package logic

import (
	"bytes"
	"comyun/dao/mysql"
	"comyun/models"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/exp/rand"
	"gopkg.in/gomail.v2"
	"html/template"
	"net/smtp"
	"os"
	"regexp"
	"time"
)

var ErrorRcpt = errors.New("recipient error")
var ErrorNotMatch = errors.New("email not match")

const EmailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

//使用第三方库
func SendEmail(email string)(code int,err error){
	// 设置SMTP服务器配置
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", viper.GetString("email.sender")) // 发件人邮箱
	mailer.SetHeader("To", email)        // 收件人邮箱
	mailer.SetHeader("Subject", "ACAT平台")
	//HTML格式
	//设置消息
	rand.Seed(uint64(time.Now().UnixNano()))
	message := 100000+rand.Intn(899999)
	////直接发送文本形式
	//mailer.SetBody("text/html", fmt.Sprintf("你的验证码为%d，有效期为十分钟，假如并非本人操作请无视",message))

	//发送HTML形式
	tmplByte,err := os.ReadFile("./logic/email.html")
	if err != nil{
		zap.L().Error("Read File email.html error"+err.Error())
		return 0,err
	}

	tmpl,err := template.New("email").Parse(string(tmplByte))
	if err != nil{
		zap.L().Error("Parse html error"+err.Error())
		return 0,err
	}

	var body bytes.Buffer
	tmpl.Execute(&body, struct {
		Code int
	}{Code: message})
	mailer.SetBody("text/html",body.String())

	smtpHost := viper.GetString("email.host_qq")
	smtpPort := viper.GetInt("email.port")// 465或587，根据邮箱服务商和是否使用SSL/TLS选择端口
	smtpUser := viper.GetString("email.sender")
	smtpPass := viper.GetString("email.pass_code") // 这不是你的邮箱密码，而是开启SMTP服务后获得的授权码

	// 构建SMTP客户端
	dialer := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	dialer.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
	} // 忽略证书校验，仅用于测试环境

	// 发送邮件
	if err = dialer.DialAndSend(mailer); err != nil {
		zap.L().Error("send email error"+err.Error())
		return 0,err
	}

	return message,nil
}

// EmailMatch 监测邮箱格式是否正确
func EmailMatch(email string) error {
	reg,err := regexp.Compile(EmailRegex)
	if err != nil{
		zap.L().Error("error to Compile regexp"+err.Error())
		return err
	}

	ok := reg.MatchString(email)
	if !ok {
		zap.L().Error("error to MatchString"+err.Error())
		return ErrorNotMatch
	}

	return nil
}

// EmailExist 判断邮箱存在问题是否正确
func EmailExist(ep *models.EmailParams)(exist bool){
	exist = mysql.EmailExist(ep.Email)

	//注册时需要邮箱不存在在
	if ep.Model == "signup"{
		if exist {
			return false
		}
		return true
	} else {
		if exist {
			return true
		}
		return false
	}
}

// SendEmailByOwn 原生发送邮箱，还未找出原因
func SendEmailByOwn(email string)(string,error){
	//设置QQ浏览器的SMTP服务器配置
	smtpHost := "smtp.qq.com"
	smtpPort := "465"

	//设置STMP认证信息
	auth := smtp.PlainAuth("",viper.GetString("email.sender"),
		viper.GetString("email.pass_code"),
		smtpHost)

	//建立STMP的tcp连接
	client,err := NewTLSClient(smtpHost,smtpPort,auth,email)
	if err != nil{
		zap.L().Error("NewSSLClient error"+err.Error())
		return "",err
	}
	defer client.Close()

	////设置TLS配置
	//tlsConfig := &tls.Config{
	//	InsecureSkipVerify: true,//开发环境使用
	//	ServerName: smtpHost,
	//}
	////建立基于TLS的tcp连接
	//conn,err := tls.Dial("tcp",smtpHost+":"+smtpPort,tlsConfig)
	//if err != nil{
	//	zap.L().Error("SSL Dial error"+err.Error())
	//	return "",err
	//}
	//defer conn.Close()
	//
	////建立SMTP客户端
	//client,err := smtp.NewClient(conn,smtpHost)
	//if err != nil{
	//	zap.L().Error("SSL NewClient error"+err.Error())
	//	return "",err
	//}
	//
	////// 启动TLS会话
	////if err := client.StartTLS(tlsConfig); err != nil {
	////	log.Fatalf("Starting TLS failed: %v", err)
	////}
	//
	////进行SMTP认证
	//if err = client.Auth(auth);err != nil{
	//	zap.L().Error("SMTP 认证 error"+err.Error())
	//	return "",err
	//}
	//
	////设置发送者
	//if err = client.Mail(viper.GetString("email.sender"));err != nil{
	//	zap.L().Error("set sender error "+err.Error())
	//	return "",err
	//}
	////设置接收者
	//if err = client.Rcpt(email);err != nil{
	//	zap.L().Error("set recipient error "+err.Error())
	//	return "",ErrorRcpt
	//}



	//发送信息
	w,errData := client.Data()
	if errData != nil{
		zap.L().Error("client.Data error"+err.Error())
		return "",errData
	}
	defer w.Close()

	//设置消息
	rand.Seed(uint64(time.Now().UnixNano()))
	message := 100000+rand.Intn(899999)
	msg := []byte(fmt.Sprintf(
		"To: %s\r\n" +
		"Subject: ACAT低代码平台\r\n" +
		"\r\n" +
		"你好，你的验证码为%d，请勿泄露他人，有效期10分钟，如果并未本人操作请无视\r\n",email,message))

	_,err = w.Write(msg)
	if err != nil{
		zap.L().Error("write email error"+err.Error())
		return "",err
	}

	return fmt.Sprintf("%d",message),nil
}

func NewTLSClient(smtpHost string,smtpPort string,auth smtp.Auth,email string) (client *smtp.Client,err error){
	//设置TLS配置
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,//开发环境使用
		ServerName: smtpHost,
	}
	//建立基于TLS的tcp连接
	conn,err := tls.Dial("tcp",smtpHost+":"+smtpPort,tlsConfig)
	if err != nil{
		zap.L().Error("SSL Dial error"+err.Error())
		return
	}
	defer conn.Close()

	//建立SMTP客户端
	client,err = smtp.NewClient(conn,smtpHost)
	if err != nil{
		zap.L().Error("SSL NewClient error"+err.Error())
		return
	}

	//// 启动TLS会话
	//if err := client.StartTLS(tlsConfig); err != nil {
	//	log.Fatalf("Starting TLS failed: %v", err)
	//}

	//进行SMTP认证
	if err = client.Auth(auth);err != nil{
		zap.L().Error("SMTP 认证 error"+err.Error())
		return
	}

	//设置发送者
	if err = client.Mail(viper.GetString("email.sender"));err != nil{
		zap.L().Error("set sender error "+err.Error())
		return
	}
	//设置接收者
	if err = client.Rcpt(email);err != nil{
		zap.L().Error("set recipient error "+err.Error())
		return nil,ErrorRcpt
	}

	return client,nil
}

