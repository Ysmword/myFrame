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
