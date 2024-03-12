package dto

type UserDTO struct {
	Id         int64
	UserName   string
	NickName   string
	ProfileUrl string
	Password   string
	// methodName
	MethodName string
	// if success
	Success bool
	// session id
	GoSessionId string
}
