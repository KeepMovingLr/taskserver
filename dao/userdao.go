package dao

import (
	"github.com/KeepMovingLr/taskserver/model"
	_ "github.com/go-sql-driver/mysql"
)

// Query a single user by userName
func SelectUserByUserName(userName string) (*model.User, error) {
	//db, _ := connectMysql()
	{
		userResult := new(model.User)

		querySql := "SELECT id, user_name, nick_name, profile_url, password, salt, gmt_create, gmt_modify  FROM user WHERE user_name = ?"
		if err := DB.QueryRow(querySql, userName).Scan(&userResult.Id, &userResult.UserName, &userResult.NickName, &userResult.ProfileUrl, &userResult.Password, &userResult.Salt, &userResult.CreateTime, &userResult.ModifyTime); err != nil {
			return nil, err
		}
		return userResult, nil
	}
}

// update nickName and ProfileUrl
func UpdateNickNameAndProfileUrl(userName string, nickName string, profileUrl string) (effectedRows int64, err error) {
	stmt, _ := DB.Prepare("update user set nick_name=? , profile_url=?  where user_name=?")
	exec, err := stmt.Exec(nickName, profileUrl, userName)
	if err != nil {
		return 0, err
	}
	affected, _ := exec.RowsAffected()
	return affected, nil
}

// useless
// update profileUrl
func UpdateUserProfileUrlByUsername(userName string, profileUrl string) (effectedRows int64, err error) {
	stmt, _ := DB.Prepare("update user set profile_url=? where user_name=?")
	exec, err := stmt.Exec(profileUrl, userName)
	if err != nil {
		return 0, err
	}
	affected, _ := exec.RowsAffected()
	return affected, nil
}

// useless
// update nickName
func UpdateNickName(userName string, nickName string) (effectedRows int64, err error) {
	stmt, _ := DB.Prepare("update user set nick_name=? where user_name=?")
	exec, err := stmt.Exec(nickName, userName)
	if err != nil {
		return 0, err
	}
	affected, _ := exec.RowsAffected()
	return affected, nil
}
