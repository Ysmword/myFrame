package controllers

import (
	"helloweb/common"
	"helloweb/logger"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// tokenString token值
var tokenString string

func init(){
	common.ReadConfig()
	logger.InitLogger()
}

func TestCreateToken(t *testing.T){
	jwtCustomClaims := JWTCustomClaims{
		"1",
		"admin",
		jwt.StandardClaims{
			ExpiresAt: int64(time.Now().Add(time.Second * 5).Unix()),
			Issuer:    "test",
		},
	}

	var err error
	// 传入的值不正确
	tokenString,err = CreateToken(JWTCustomClaims{})
	if err!=nil{
		t.Log(err)
	}else{
		t.Error(tokenString)
	}

	// 传递正确的值
	tokenString,err = CreateToken(jwtCustomClaims)
	if err!=nil{
		t.Error(err)
	}else{
		t.Log(tokenString)
	}
}


func TestParseToken(t *testing.T){
	t.Log(tokenString)
	// 正确值
	jwtCustomClaims,err := ParseToken(tokenString)
	if err!=nil{
		t.Error(err)
	}else{
		t.Log(jwtCustomClaims)
	}

	// 空值
	jwtCustomClaims,err = ParseToken("")
	if err!=nil{
		t.Log(err)
	}else{
		t.Error(jwtCustomClaims)
	}

	// 时间过期
	time.Sleep(7*time.Second)
	jwtCustomClaims,err = ParseToken(tokenString)
	if err!=nil{
		t.Log(err)
	}else{
		t.Error(jwtCustomClaims)
	}
}