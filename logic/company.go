package logic

import (
	"comyun/dao/mysql"
	"comyun/models"
	"comyun/pkg/snowflake"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
	"mime/multipart"
	"strconv"
)

var ValueString = "(?,?,?,?,?,?,?,?,?,?)"

func NewCompany(p *models.Company) error {
	//在mysql中进行注册
	//生成company_id
	p.CompanyId = snowflake.GenID()
	err := mysql.NewCompany(p)
	if err != nil{
		zap.L().Error("NewCompany err"+err.Error())
		return err
	}

	//在mysql中user表对对象的所属公司进行更新
	err = mysql.UserChangeCompany(p)

	return err

}

func InviteGroup(cid string,file multipart.File) error {
	f,err := excelize.OpenReader(file)
	if err != nil{
		zap.L().Error("OpenReader err"+err.Error())
		return err
	}

	sheet,err := f.GetRows(f.GetSheetList()[0])
	if err != nil{
		zap.L().Error("getRows err"+err.Error())
		return err
	}

	contents := make([]interface{},0,(len(sheet)-2)*10)
	valueStrings := make([]string,0,len(sheet)-2)
	for i,row := range sheet{
		if i < 2 {
			continue
		}
		valueStrings = append(valueStrings,ValueString)
		contents = append(contents,snowflake.GenID())
		contents = append(contents,false)
		contents = append(contents,cid)
		contents = append(contents,row[0])
		gender,err := strconv.ParseBool(row[1])
		if err != nil{
			return err
		}
		contents = append(contents,gender)
		contents = append(contents,row[2])
		contents = append(contents,row[3])
		contents = append(contents,row[4])
		contents = append(contents,row[5])
		contents = append(contents,row[6])
	}

	err = mysql.UserXlsxGroup(contents,valueStrings)
	if err != nil{
		zap.L().Error("UserXlsxGroup err"+err.Error())
		return err
	}

	return nil
}
