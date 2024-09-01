package mysql

import "comyun/models"

func GetLoggerPage(cl *models.ComLogParam)([]models.ComLogMess,error){
	sqlstr := `select * from logger where company_id = ? limit ? offset ?`

	mess := new([]models.ComLogMess)
	if err := db.Select(mess,sqlstr,cl.CompanyID,cl.Sum,(cl.Page-1)*cl.Sum);err != nil{
		return nil,err
	}
	return *mess,nil
}

func GetLoggerPerson(cl *models.ComLogParam)([]models.ComLogMess,error){
	sqlstr := `select * from logger 
		where company_id = ? AND user_id = ? 
		limit ? offset ?`

	mess := new([]models.ComLogMess)
	if err := db.Select(mess,sqlstr,cl.CompanyID,cl.UserId,cl.Sum,(cl.Page-1)*cl.Sum);err != nil{
		return nil,err
	}
	return *mess,nil
}

func GetLoggerTime(cl *models.ComLogParam)([]models.ComLogMess,error){
	sqlstr := `select * from logger 
		where company_id = ? AND log_time BETWEEN ? AND ? 
		limit ? offset ?`

	mess := new([]models.ComLogMess)
	if err := db.Select(mess,sqlstr,cl.CompanyID, cl.StartTime,cl.EndTime,cl.Sum,(cl.Page-1)*cl.Sum);err != nil{
		return nil,err
	}
	return *mess,nil
}

func GetLoggerPT(cl *models.ComLogParam)([]models.ComLogMess,error){
	sqlstr := `select * from logger 
		where company_id = ? AND user_id = ? AND log_time BETWEEN ? AND ? 
		limit ? offset ?`

	mess := new([]models.ComLogMess)
	if err := db.Select(mess,sqlstr,cl.CompanyID, cl.UserId,cl.StartTime,cl.EndTime,cl.Sum,(cl.Page-1)*cl.Sum);err != nil{
		return nil,err
	}
	return *mess,nil
}

func SetLog(log *models.Log) error {
	sqlstr := `insert into logger(
		user_id,user_name,company_id,department,role,log_time) 
		values (:user_id,:user_name,:company_id,:department,:role,:log_time)`

	//命名传输
	_,err := db.NamedExec(sqlstr,log)
	return err
}
