package coreservice

import (
	"fmt"
	redisconn "github.com/KeepMovingLr/taskserver/cache"
	"github.com/KeepMovingLr/taskserver/constant"
	"github.com/KeepMovingLr/taskserver/dao"
	"github.com/KeepMovingLr/taskserver/dto"
	"github.com/KeepMovingLr/taskserver/utils"
	"strings"
)

// check weather a user can login. if success, return the user info
func LoginAuthenticate(userName, password string) (*dto.UserDTO, error) {

	// input parameter check
	if userName == "" || password == "" {
		return nil, fmt.Errorf("userName or password is blank")
	}

	cache := redisconn.MyCache.GetFromCache(userName)

	result := new(dto.UserDTO)
	result.UserName = userName
	if cache != nil {
		userString := cache.(string)
		split := strings.Split(userString, constant.CacheSplit)

		result.NickName = string(split[1])
		result.ProfileUrl = string(split[2])
		result.Success = true
		if !utils.CheckSha256Password(password, string(split[0])) {
			return nil, fmt.Errorf("password not correct")
		}
		return result, nil
	} else {
		user, err := dao.SelectUserByUserName(userName)
		if err != nil {
			return nil, fmt.Errorf("dao operation false,%v", err)
		}
		if user == nil {
			return nil, fmt.Errorf("user not exist")
		}
		pwdFromStorage := user.Password
		result.NickName = user.NickName
		result.ProfileUrl = user.ProfileUrl
		result.Success = true
		if utils.CheckSha256Password(password, string(pwdFromStorage)) {
			redisconn.MyCache.AddToCache(userName, user.Password+constant.CacheSplit+user.NickName+constant.CacheSplit+user.ProfileUrl)
			return result, nil
		} else {
			return nil, fmt.Errorf("password not correct")
		}

	}
}

// modify profileUrl by username
func ModifyUserInfoByUserName(userName string, nickName, profileUrl, sessionId string) (*dto.UserDTO, error) {
	// input parameter check
	if userName == "" {
		return nil, fmt.Errorf("userName is blank")
	}
	// login status check
	userNameFromSession, _ := redisconn.Get(sessionId)
	if userNameFromSession == nil {
		// 未登录
		return nil, fmt.Errorf("userName not login")
	}
	if string(userNameFromSession) != userName {
		return nil, fmt.Errorf("authenticate failed")
	}
	effectedRow, err := dao.UpdateNickNameAndProfileUrl(userName, nickName, profileUrl)
	if err != nil {
		return nil, fmt.Errorf("dao operation false,%v", err)
	}
	userResult := dto.UserDTO{
		Success:     false,
		UserName:    userName,
		GoSessionId: sessionId,
	}
	if effectedRow == 0 {
		return &userResult, nil
	}
	// invalid cache
	redisconn.MyCache.InvalidCacheByKey(userName)
	// return result
	userResult.Success = true
	userResult.NickName = nickName
	userResult.ProfileUrl = profileUrl
	return &userResult, nil
}

// modify nickName by username
func ModifyNickNameByUsername(userName, nickName string) (*dto.UserDTO, error) {
	effectedRow, err := dao.UpdateNickName(userName, nickName)
	if err != nil {
		return nil, fmt.Errorf("dao operation false,%v", err)
	}
	userResult := dto.UserDTO{
		Success:  false,
		UserName: userName,
		NickName: nickName,
	}
	if effectedRow == 0 {
		return &userResult, nil
	}
	userResult.Success = true
	return &userResult, nil
}
