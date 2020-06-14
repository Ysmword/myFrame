package controllers

import (
	"fmt"
	"helloweb/common"
	"helloweb/logger"

	"errors"

	"github.com/dgrijalva/jwt-go"
)

// JWTCustomClaims 将生成token需要的信息，封装到一个结构体中
type JWTCustomClaims struct{
	// 与前端需要数据交互信息,这是只是设置一个大概,但是大致的思路为：权限是以角色为基础分配的
	// UID 用户id
	UID string `json:"uid"`
	Role string `json:"role"`

	jwt.StandardClaims `json:"jwt"`
}

// ErrInvalidToken token是否有效
var ErrInvalidToken = errors.New("token is invalid")

// CreateToken 生成token
func CreateToken(jwtCustomClaims JWTCustomClaims)(tokenString string,err error){

	// if jwtCustomClaims==nil{
	// 	err = fmt.Errorf("CreateToken jwtCustomClaims is nil")
	// 	logger.Z.Error(err.Error())
	// 	return "",err
	// }

	if jwtCustomClaims.UID==""{
		err = fmt.Errorf("CreateToken jwtCustomClaims.UID is null")
		logger.Z.Error(err.Error())
		return "",err
	}

	if jwtCustomClaims.Role==""{
		err = fmt.Errorf("CreateToken jwtCustomClaims.Role is null")
		logger.Z.Error(err.Error())
		return "",err
	}


	if common.Conf.SecretKey.SecretKey == "" {
		err = fmt.Errorf("common.Conf.SecretKey.SecretKey is null")
		logger.Z.Error(err.Error())
		return "",err
	}

	// 生成token
	logger.Z.Info("生成token " + common.Conf.SecretKey.SecretKey)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtCustomClaims)
	tokenString, err = token.SignedString([]byte(common.Conf.SecretKey.SecretKey))
	if err!=nil{
		logger.Z.Error(err.Error())
		return "",err
	}
	logger.Z.Info(tokenString)
	return
}

// ParseToken 解析token
func ParseToken(tokenString string)(*JWTCustomClaims,error){

	if tokenString == ""{
		err := fmt.Errorf("ParseToken tokenString is null")
		logger.Z.Error(err.Error())
		return nil,err
	}

	if common.Conf.SecretKey.SecretKey == "" {
		err := fmt.Errorf("common.Conf.SecretKey.SecretKey is null")
		logger.Z.Error(err.Error())
		return nil,err
	}

	token,err := jwt.ParseWithClaims(tokenString,&JWTCustomClaims{},func(token *jwt.Token) (interface{}, error) {
		return []byte(common.Conf.SecretKey.SecretKey), nil
	})
	if err!=nil{
		logger.Z.Error(err.Error())
		return nil,err
	}

	if !token.Valid{
		logger.Z.Error(ErrInvalidToken.Error())
		return nil,ErrInvalidToken
	}
	logger.Z.Info("解析token")
	jwtCustomClaims,ok := token.Claims.(*JWTCustomClaims)
	if !ok{
		err = fmt.Errorf("token.Claims can not convert to *JWTCustomClaims")
		logger.Z.Error(err.Error())
		return nil,err
	}
	
	logger.Z.Info(jwtCustomClaims.Role)

	return jwtCustomClaims,nil
}





