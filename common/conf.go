package common

import (

	"log"

	"gopkg.in/gcfg.v1"
)

// Config 配置
type Config struct {
	// AppConf app基础配置
	AppConf struct {
		// httpport 监听端口
		Httpport string
	}
	// Addr 服务器的地址
	Addr struct { 
		IP string
	}
	// Postgres 数据库配置
	Postgres struct {
		// Host 主机
		Host string
		// Port 端口号
		Port int
		// Sslmode 模式
		Sslmode string
		// User 使用者
		User string
		// Password 密码
		Password string
		// Dbname 数据库名字
		Dbname string
	}
	// Git git 所需要的数据
	Git struct{
		// GitEmail 用户邮箱
		GitEmail string
		// GitName 用户账号
		GitName string
		// GitPW 用户密码
		GitPW string
	}
}

// Conf 配置变量
var Conf = new(Config)

// ReadConfig 读取配置文件
func ReadConfig()  {
	if HomePath == "" {
		log.Fatal("初始化app路径失败，HomePath = ", HomePath)
	}

	configPath := HomePath + "/" + "conf.ini"

	log.Println("配置文件路径",configPath)

	err := gcfg.ReadFileInto(Conf,configPath)
	if err!=nil{
		log.Fatal("gcfg.ReadFileInto error:",err)
	}

	log.Printf(Conf.AppConf.Httpport)
}

