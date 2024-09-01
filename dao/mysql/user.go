package mysql

import (
	"comyun/models"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"strings"
)

var ErrorPassword = errors.New("error to password")

func SignupUser(u *models.User) error {
	sqlstr := `insert into user(
		user_id,boss,company_id,name,gender,telephone,email,password,department,role) 
		values (?,?,?,?,?,?,?,?,?,?)`

	_,err := db.Exec(sqlstr,u.UserID,u.Boss,u.CompanyID,u.Name,u.Gender,u.Telephone,u.Email,u.Password,u.Department,u.Role)
	if err != nil{
		return err
	}

	return nil
}

func GetUserID(email string)(int64,error){
	sqlstr := "select user_id from user where email = ?"

	var id int64
	err := db.Get(&id,sqlstr,email)
	if err != nil{
		zap.L().Error(err.Error())
	}

	return id,err
}

func LoginPass(up *models.UserParams)(int64,error){
	sqlstr := "select id,password from user where telephone = ?"

	u := &struct {
		Id int64`json:"id"`
		PassWord string`json:"pass_word"`
	}{}
	err := db.Get(u,sqlstr,up.Telephone)

	if err != nil{
		zap.L().Error("Get id,password err"+err.Error())
		return 0,err
	}
	if u.PassWord != up.Password{
		return 0,ErrorPassword
	}

	return u.Id,nil
}

// 更新员工的公司信息
func UserChangeCompany(p *models.Company) error {
	sqlStr := `update user set 
		company_id = ?,boss = ?
		where user_id = ?`
	_,err := db.Exec(sqlStr,p.CompanyId,1,p.OwnerId)
	//返回错误
	return err
}

func GetUserFromEmail(email string) (*models.Log,error) {
	sqlstr := `select user_id,user_name,company_id,department,role from user where email = ?`

	log := new(models.Log)
	err := db.Select(log,sqlstr,email)

	return log,err
}

func GetUserFromTele(tele string) (*models.Log,error) {
	sqlstr := `select user_id,user_name,company_id,department,role from user where telephone = ?`

	log := new(models.Log)
	err := db.Select(log,sqlstr,tele)

	return log,err
}

// UserXlsxGroup 实现用户的批量导入
func UserXlsxGroup(contents []interface{},valueString []string) error {
	//fmt.Println(contents)
	sqlstr := fmt.Sprintf(`insert into user 
		(user_id,boss,company_id,name,gender,telephone,email,password,department,role) 
		values %s`,strings.Join(valueString,","))

	//ns,err := db.PrepareNamed(sqlstr)
	//if err != nil{
	//	zap.L().Error("db.PrepareNamed error to "+err.Error())
	//	return err
	//}
	_,err := db.Exec(sqlstr,contents...)
	if err != nil{
		zap.L().Error("db.Exec error to "+err.Error())
		return err
	}

	return err
}

func GetCompanyID(userId int64) int64 {
	sqlstr := "select company_id from user where user_id = ?"

	var id int64
	if err := db.Get(&id,sqlstr,userId);err != nil{
		zap.L().Error("select company_id error"+err.Error())
		return 0
	}

	return id
}

func GetRole(userID int64) (int,error) {
	sqlstr := "select role from user where user_id = ?"

	role := 0
	err := db.Get(&role,sqlstr,userID)
	if err != nil{
		zap.L().Error("error select role"+err.Error())
	}

	return role,err
}
