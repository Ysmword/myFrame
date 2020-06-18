package controllers

import (
	"fmt"
	"helloweb/logger"
	"net/http"
	"strings"
)


// Hello 哈喽世界
func Hello(w http.ResponseWriter,r *http.Request)(interface{},error){

	if strings.ToLower(r.Method)!="get"{
		err := fmt.Errorf("请求方式必须是get")
		logger.Z.Info(err.Error())
		return nil,err
	}

	return "hello world",nil
}