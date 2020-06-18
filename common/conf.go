package common

import (
	"log"

	"gopkg.in/gcfg.v1"

	valid "github.com/asaskevich/govalidator"
)

// ConfigPath 配置文件的路径
var ConfigPath string

// Config 配置
type Config struct {
	AppConf `valid:"required"`
	Addr    `valid:"required"`
	Postgres
	Git
	SecretKey
	HotUpate
	ExecutablePath
}

// AppConf app基础配置
type AppConf struct {
	// httpport 监听端口
	Httpport string `json:"httpport" valid:"required,numeric"` // required表示一定要进行数据验证，numeric表示是数字
}

// Addr 服务器的地址
type Addr struct {
	IP string `json:"ip" valid:"required,ip"`
}

// Postgres 数据库配置
type Postgres struct {
	// Host 主机
	Host string `json:"host" valid:"ip,optional"`
	// Port 端口号
	Port int `json:"port" valid:"int,optional"`
	// Sslmode 模式
	Sslmode string `json:"sslmode" valid:"sslmode,optional"`
	// User 使用者
	User string `json:"user" valid:"type(string)"`
	// Password 密码
	Password string `json:"password" valid:"type(string)"`
	// Dbname 数据库名字
	Dbname string `json:"dbname" valid:"type(string)"`
}

// Git 所需要的数据
type Git struct {
	// GitEmail 用户邮箱"
	GitEmail string `json:"GitEmail" valid:"email,optional"`
	// GitName 用户账号
	GitName string `json:"gitName" valid:"type(string)"`
	// GitPW 用户密码
	GitPW string `json:"gitpw" valid:"type(string)"`
}

// SecretKey 密钥
type SecretKey struct {
	SecretKey string `json:"secretKey " valid:"type(string)"`
}

// HotUpate 热更新
type HotUpate struct {
	IsOpen bool `json:"isOpen"`
}

// ExecutablePath 可执行文件路径
type ExecutablePath struct {
	WinExecutablePath   string `json:"WinExecutablePath"`
	LinuxExecutablePath string `json:"LinuxExecutablePath"`
}

func init() {

	// valid.SetFieldsRequiredByDefault(true)
	// 自定义添加验证器，要放在函数中
	valid.TagMap["sslmode"] = valid.Validator(func(sslmode string) bool {
		var sslmodes = []string{
			"disable",
			"allow",
			"prefer",
			"require",
			"verify-ca",
			"verify-full",
		}
		for i := 0; i < len(sslmodes); i++ {
			if sslmodes[i] == sslmode {
				return true
			}
		}

		log.Println("sslmode value is not disable or allow ro prefer or require or verify-ca or verify-full")
		return false
	})
	if HomePath == "" {
		log.Fatal("初始化app路径失败，HomePath = ", HomePath)
	}
	ConfigPath = HomePath + "/" + "conf.ini"
}

// Conf 配置变量
var Conf = new(Config)

// ReadConfig 读取配置文件
func ReadConfig() {

	log.Println("配置文件路径", ConfigPath)

	err := gcfg.ReadFileInto(Conf, ConfigPath)
	if err != nil {
		log.Fatal("gcfg.ReadFileInto error:", err)
	}

	result, err := valid.ValidateStruct(Conf)

	if err != nil {
		log.Fatal("error:", err)
	}

	log.Println("验证配置文件结构的有效性", result)

	log.Printf(Conf.AppConf.Httpport)
}
