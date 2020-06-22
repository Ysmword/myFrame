package controllers

import (
	"fmt"
	"helloweb/logger"
	"net/http"
	"strings"
	"time"
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


// Hello1 哈喽世界
func Hello1(w http.ResponseWriter,r *http.Request)(interface{},error){

	if strings.ToLower(r.Method)!="get"{
		err := fmt.Errorf("请求方式必须是get")
		logger.Z.Info(err.Error())
		return nil,err
	}
	logger.Z.Info("hello world1")
	return "hello world1",nil
}


// Hello2 哈喽世界
func Hello2(w http.ResponseWriter,r *http.Request)(interface{},error){

	if strings.ToLower(r.Method)!="get"{
		err := fmt.Errorf("请求方式必须是get")
		logger.Z.Info(err.Error())
		return nil,err
	}
	logger.Z.Info("hello world2")
	time.Sleep(2*time.Second)
	return "hello world2",nil
}