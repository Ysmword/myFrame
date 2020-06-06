package controllers

import (
	"net/http"
)

// StaticServe 文件服务器 http.Handle("/file",http.StripPrefix("/file/", http.FileServer(http.Dir("."))))
func StaticServe(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	http.StripPrefix("/file/", http.FileServer(http.Dir("."))).ServeHTTP(w, r)
	return nil, nil
}
