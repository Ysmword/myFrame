package common

import (
	"log"
	"os"
	"path/filepath"

	"github.com/gorilla/websocket"
	"gopkg.in/gomail.v2"
)

var (
	// HomePath application startup directory
	HomePath string
	// WsUserTable 创建一个websocket用户表
	WsUserTable = make(map[*websocket.Conn]string)

	// LastModifiedTime 最新修改时间
	LastModifiedTime string
	

	// 发送用户
	senderAddr    = "ysm_515@163.com"
	// 接收用户组
	receviceAddrs = []string{
		"1843121593@qq.com",
	}
	// 授权码,用于发送邮件信息
	authCode = "ysm121388"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	var err error
	HomePath, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	// 查看整个项目文件是否有增加
	fileInfo,err :=os.Stat(HomePath)
	if err!=nil{
		log.Println(err)
		return
	}
	// 记录修改时间
	LastModifiedTime = fileInfo.ModTime().String()
}

// SendMail 使用框架发送邮件（不包含文件），毕竟是使用第三方的东西，学校网，不能够发送
func SendMail(msgInfo string) error {

	sendData := msgInfo + "<br>"
	// 创建消息体，不要认为Header完全看作是请求头，Header包含了邮件的头部信息
	m := gomail.NewMessage()
	// 主题
	subject := "app错误信息"
	// // 匿名
	// nickName := "nickName"
	// 目标
	dest := []string{}
	m.SetHeader("From", m.FormatAddress(senderAddr, "nickname"))
	for _, value := range receviceAddrs {
		dst := m.FormatAddress(value, "nickname")
		dest = append(dest, dst)
	}
	dst := m.FormatAddress(senderAddr, "nickname") // 记得也要发送一个内容给自己。以证明自己的清白
	dest = append(dest, dst)
	m.SetHeader("To", dest...)
	m.SetHeader("Subject", subject)
	// "Content-Type: text/html; charset=UTF-8"
	m.SetBody("text/html", sendData)                                 //设置邮件发送的内容
	d := gomail.NewDialer("smtp.163.com", 465, senderAddr, authCode) // 创建连接
	if err := d.DialAndSend(m); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// IsDir 判断目录是否存在
func IsDir(fileAddr string) (bool, error) {
	s, err := os.Stat(fileAddr)
	if err != nil {
		return false, err
	}
	return s.IsDir(), nil
}

// IsExist 判断文件是否存在
func IsExist(fileAddr string) bool {
	// 读取文件信息，判断文件是否存在
	_, err := os.Stat(fileAddr)
	if err != nil {
		if os.IsExist(err) { // 根据错误类型进行判断
			return true
		}
		return false
	}
	return true
}
