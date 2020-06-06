package controllers

import (
	"fmt"
	"net/http"
	"helloweb/logger"
	"helloweb/common"
	"io/ioutil"
	"encoding/base64"
)

// StaticServe 文件服务器 http.Handle("/file",http.StripPrefix("/file/", http.FileServer(http.Dir("."))))
func StaticServe(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	http.StripPrefix("/file/", http.FileServer(http.Dir("."))).ServeHTTP(w, r)
	return nil, nil
}

// Save 保存图片
func Save(rootPath string,nowTimeStr string,data string) error {
	if rootPath==""{
		err := fmt.Errorf("Save rootPath is null")
		logger.Z.Error(err.Error())
		return err
	}
	if nowTimeStr==""{
		err := fmt.Errorf("Save nowTimeStr is null")
		logger.Z.Error(err.Error())
		return err
	}
	if data==""{
		err := fmt.Errorf("Save data is null")
		logger.Z.Error(err.Error())
		return err
	}
	isExist, err := common.IsDir(rootPath)
	if err != nil {
		logger.Z.Error(err.Error())
		return err
	}
	if isExist {
		savePath := rootPath + nowTimeStr + ".png"
		if !common.IsExist(savePath) {
			buf, err := base64.StdEncoding.DecodeString(data)
			if err != nil {
				logger.Z.Error(err.Error())
				return err
			}
			logger.Z.Info("写入文件")
			err = ioutil.WriteFile(savePath, buf, 0666)
			if err != nil {
				logger.Z.Error(err.Error())
				return  err
			}
		}
	}
	return nil
}
