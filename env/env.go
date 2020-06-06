package env

import(
	"log"
	"fmt"
	"time"
	"helloweb/conf"
	"helloweb/logger"
	"helloweb/models"
	"helloweb/routers"
)

// InitAll 初始化项目  fdsafsd
func InitAll(){

	log.Println("读取配置文件")
	conf.ReadConfig()

	log.Println("初始化日志文件")
	logger.InitLogger()

	log.Println("初始化数据库")
	models.InitDB()

	log.Println("开启服务")
	s := `
	***************************************************************
	******** app start: %s
	***************************************************************`
	logger.Z.Info(fmt.Sprintf(s, time.Now().String()))
	routers.StartServer()
}