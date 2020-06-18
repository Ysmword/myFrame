package routers

import (
	"context"
	"fmt"
	"helloweb/common"
	"helloweb/controllers"
	"helloweb/logger"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// mux 服务对象表
var serviceObjectTable = make(map[string]*ControllerInfo)

// hellowebHTTPHandler 广州美术学院同学的app
type hellowebHTTPHandler struct{}

// specialString 不必要的路由
var specialString = "/favicon.ico"

// LastModifiedTime 最新修改时间
var LastModifiedTime string

// ControllerInfo 保存有关控制器的信息
type ControllerInfo struct {
	// url path
	Path string

	// 处理函数
	Fn func(http.ResponseWriter, *http.Request) (interface{}, error)

	// 作用介绍
	APIName string

	// 该控制器是否可用,用于调试
	Available bool

	// 是否是文件服务器
	isFileSystem bool

	// 是否是websocket接口
	isWebSocket bool
}

func init() {
	fileInfo, err := os.Stat(common.ConfigPath)
	if err != nil {
		log.Fatal(err)
	}
	LastModifiedTime = fileInfo.ModTime().String()
	log.Println("LastModifiedTime", LastModifiedTime)
}

// ServiceHttp 服务
func (g *hellowebHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if strings.Contains(r.URL.String(), specialString) {
		logger.Z.Warn("防止angular发送诡异的请求")
		return
	}

	logger.Z.Info("开始访问")

	// 注意有get请求或者是delete请求的部分
	method := strings.ToLower(r.Method)
	logger.Z.Info("method = " + method)
	delMethod := "delete"  // delete请求方式
	getMethod := "get"     // get请求方式
	api := ""              // api，用来接收接口具体名字的
	var err error          // 错误信息
	var result interface{} // 接收信息
	if method == delMethod {
		api = strings.Split(r.URL.String(), "?")[0]
	} else if method == getMethod {
		// 判断是不是文件文件服务器
		if strings.Contains(r.URL.String(), "file") {
			api = "/file"
		} else {
			api = strings.Split(r.URL.String(), "?")[0]
		}
	} else {
		api = r.URL.String()
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Action, Module")
	logger.Z.Info(api)
	// h找到对应的对象，ok表示一种状态
	if h, ok := serviceObjectTable[api]; ok {
		// 一开始先判断是不是静态文件服务器接口
		if h.isFileSystem || h.isWebSocket {
			// 文件服务器
			h.Fn(w, r)
		} else {
			result, err = h.Fn(w, r)
			if err != nil {
				err = fmt.Errorf("运行接口%v的时候报错，err:%v", r.URL.String(), err)
				logger.Z.Error(err.Error())
				if err = controllers.ErrorResp(w, r, err.Error()); err != nil {
					logger.Z.Error(r.URL.String() + ",运行报错")
				}
			} else {
				if h.Available == true {
					if err = controllers.SuccessResp(w, r, result); err != nil {
						logger.Z.Error(r.URL.String() + ",运行报错")
					}
				} else {
					logger.Z.Info("该接口不可用")
					if err = controllers.ErrorResp(w, r, "该接口不可用"); err != nil {
						logger.Z.Error(r.URL.String() + ",运行报错")
					}
				}
			}
		}
	} else {
		logger.Z.Warn("mux[r.URL.String()]转化失败")
		if err = controllers.ErrorResp(w, r, "mux[r.URL.String()]转化失败"); err != nil {
			logger.Z.Info(r.URL.String() + ",运行报错")
		}
	}
	if err != nil {
		logger.Z.Info("发送错误信息")
		err = common.SendMail(err.Error())
		if err != nil {
			logger.Z.Error(err.Error())
		}
	}
}

// StartServer 开启服务 GracefulFlag 是否优雅重启
func StartServer(GracefulFlag bool) {
	
	// 这里进行路由注册 serviceObjectTable["/exmple"] = &ControllerInfo{....}
	serviceObjectTable["/hello"] = &ControllerInfo{Path: "/hello",Fn: controllers.Hello,APIName:"哈喽世界",Available:true,isFileSystem:false,isWebSocket:false}
	// 这里搞个可以搞一个热更新
	if !common.Conf.HotUpate.IsOpen {
		server := http.Server{
			Addr:    ":" + common.Conf.AppConf.Httpport,
			Handler: &hellowebHTTPHandler{},
		}
		if err := server.ListenAndServe(); err != nil {
			log.Fatal("开启服务失败", err)
		}
	} else {
		var err error
		var listener net.Listener
		if  GracefulFlag {
			f := os.NewFile(3, "")
			listener, err = net.FileListener(f)
			log.Println("优雅执行热更新")
		} else {
			listener, err = net.Listen("tcp", ":"+common.Conf.AppConf.Httpport)
			if err != nil {
				log.Fatal(err)
			}
		}
		server := http.Server{
			Handler: &hellowebHTTPHandler{},
		}

		go func() {
			err := server.Serve(listener)
			if err != nil {
				log.Fatal(err)
			}
		}()
		ListenHandler(&server, listener)
	}
}

// ListenHandler 信号处理器，这里监听配置文件的变化
func ListenHandler(server *http.Server, listener net.Listener) {

	// 设置十秒的的监听器
	ticker := time.NewTicker(10 * time.Second)

	for {
		select {
		case <-ticker.C:
			var err error
			var fileInfo os.FileInfo
			var cmd *exec.Cmd
			fileInfo, err = os.Stat(common.ConfigPath)
			if err != nil {
				log.Println("重启不成功")
				log.Println("错误原因为：", err)
				break
			}
			if fileInfo.ModTime().String() != LastModifiedTime {
				LastModifiedTime = fileInfo.ModTime().String()
				// 这里优雅的关闭已有的服务，防止没有运行完以后的连接
				ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
				// 父子进程拷贝
				t1, ok := listener.(*net.TCPListener)
				if !ok {
					log.Println("listener is not tcp listener")
					log.Println("重启不成功")
					break
				}
				currentFD, err := t1.File()
				if err != nil {
					log.Println("acquiring listener file failed", err)
					log.Println("重启不成功")
					break
				}
				if runtime.GOOS == "windows" {
					cmd = exec.Command(common.Conf.ExecutablePath.WinExecutablePath, "-g")
				} else if runtime.GOOS == "linux" {
					cmd = exec.Command("./"+common.Conf.ExecutablePath.LinuxExecutablePath, "-g")
				}
				log.Println(cmd.Args)
				cmd.ExtraFiles = []*os.File{currentFD}
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				err = cmd.Start()
				if err != nil {
					log.Println("重启不成功")
					log.Println("错误原因为：", err)
					break
				}
				// 优雅的服务,放在前面也不行，因为会立即关闭的
				server.Shutdown(ctx)
			}

		}
	}

}
