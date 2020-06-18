package env

import (
	"flag"
	"fmt"
	"helloweb/common"
	"helloweb/logger"
	"helloweb/models"
	"helloweb/routers"
	"log"
	"time"
	"os"
	// valid "github.com/asaskevich/govalidator"
)

// GitFlag 是否开启自动提交，默认值not
var GitFlag = flag.String("acp", "not", "enable auto submit yes or not,defualt value is not")

// InitAll 初始化项目
func InitAll() {

	common.Monitor()

	log.Println("读取配置文件")
	common.ReadConfig()

	flag.Parse()
	switch *GitFlag {
	case "not":
		log.Println("不开启自动提交代码到仓库的功能")
	case "yes":
		log.Println("开启自动提交代码到仓库的功能")
		err := common.GitCmd()
		if err != nil {
			log.Println(err)
		}
		os.Exit(0)
	default:
		log.Println("输入不正确，可以参考：-acp=yes")
		os.Exit(0)
	}

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
