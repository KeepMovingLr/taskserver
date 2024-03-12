package model

type User struct {
	Id         int64
	UserName   string
	NickName   string
	ProfileUrl string
	Password   string
	Salt       string
	CreateTime int64
	ModifyTime int64
}
