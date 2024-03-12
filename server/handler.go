package server

import (
	"github.com/KeepMovingLr/taskserver/constant"
	"github.com/KeepMovingLr/taskserver/coreservice"
	"github.com/KeepMovingLr/taskserver/dto"
)

var HandlerMap = make(map[string]HandMethod)

// init handler map in init() method
func init() {
	HandlerMap[constant.MethodLoginCheck] = UserLoginHandler
	HandlerMap[constant.MethodUpdateUserInfo] = UserUpdateHandler
}

type HandMethod func(dto.UserDTO) dto.UserDTO

func UserLoginHandler(requestUser dto.UserDTO) (handleResult dto.UserDTO) {
	resultUser, err := coreservice.LoginAuthenticate(requestUser.UserName, requestUser.Password)
	if resultUser == nil {
		resultUser = &requestUser
		resultUser.Success = false
		return *resultUser
	}
	if err != nil {
		resultUser.Success = false
		return *resultUser
	}
	resultUser.Success = true
	return *resultUser
}

func UserUpdateHandler(requestUser dto.UserDTO) (handleResult dto.UserDTO) {
	resultUser, _ := coreservice.ModifyUserInfoByUserName(requestUser.UserName, requestUser.NickName, requestUser.ProfileUrl, requestUser.GoSessionId)
	if resultUser == nil {
		requestUser.Success = false
		return requestUser
	}
	resultUser.Success = true
	return *resultUser
}

func HandlerDispatcher(data dto.UserDTO, handler HandMethod) (handleResult dto.UserDTO) {
	return handler(data)
}
