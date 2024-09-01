package mysql

import "comyun/models"

func NewCompany(p *models.Company) error {
	sqlstr := `insert into company(
		company_id,name,owner_id) 
		values (?,?,?)`

	_,err := db.Exec(sqlstr,p.CompanyId,p.Name,p.OwnerId)

	return err
}

func InviteLogin(i *models.UserInvite) (err error) {
	sqlstr := `update user set 
		company_id = ?,boss = ?,department = ?,role = ? 
		where user_id = ?`

	_,err = db.Exec(sqlstr,i.CompanyID,0,i.Department,i.Role)

	return err
}

func GetUsers(cu *models.ComUserParam) ([]*models.User,error){
	sqlstr := `select * from user where company_id = ? limit ? offset ?`

	users := new([]*models.User)
	err := db.Select(users,sqlstr,cu.CompanyID,cu.Sum,(cu.Page-1)*cu.Sum)
	if err != nil{
		return nil,err
	}
	return *users,nil
}

