package common

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"syscall"

	// "strconv"
	"sync"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// 远程开发调试
// 生成最新编译文件
// 将编译文件上传到远程服务器
// 在远程服务器上运行

// 协程组锁
var wg sync.WaitGroup
var lock sync.Mutex

// Connect 创建sftp连接
func Connect() (*ssh.Client, *sftp.Client, error) {

	// 获取验证方式
	auth := make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(Conf.FTP.Password))

	// 获取验证配置文件
	clientConfig := &ssh.ClientConfig{
		User:    Conf.FTP.User,
		Auth:    auth,
		Timeout: 30 * time.Second, // 连接超时时间
		// 验证服务器的主机密钥 HostKeyCallback:ssh.InsecureIgnoreHostKey(),也可以
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	addr := fmt.Sprintf("%s:%d", Conf.FTP.Host, Conf.FTP.Port)

	// 创建链接
	sshClient, err := ssh.Dial("tcp", addr, clientConfig)
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}

	// session, err := sshClient.NewSession()
	// if err != nil {
	// 	log.Println(err)
	// 	return nil, nil, err
	// }

	return sshClient, sftpClient, nil
}

// UploadRun 上传可运行文件和运行文件
func UploadRun() error {
	sshClient, sftpClient, err := Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer sftpClient.Close()
	defer sshClient.Close()
	go func() {
		// 先生成可执行文件
		cmdShell := "set CGO_ENABLED=0&&set GOOS=linux&&set GOARCH=amd64&&cd " + HomePath + "&& go build"
		log.Println("cmdShell =", cmdShell)
		cmd := exec.Command("cmd", "/C", cmdShell)
		data, err := cmd.Output()
		if err != nil {
			log.Fatal(err)
		}

		if len(data) != 0 {
			log.Fatal(string(data))
		}
		// 上传可执行文件
		log.Println("上传可执行文件")
		localFilePath := HomePath + "/helloweb"
		if !IsExist(localFilePath) {
			log.Fatal(HomePath + "/helloweb" + "文件不存在")
		}
		remoteFileName, err := uploadFile(localFilePath, sftpClient)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("上传成功")
		log.Println("上传配置文件")
		_, err = uploadFile(ConfigPath, sftpClient)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("上传成功")
		// 运行命令
		session, err := sshClient.NewSession()
		if err != nil {
			log.Fatal(err)
		}
		session.Stderr = os.Stderr
		session.Stdout = os.Stdout
		shellCmd := "cd " + Conf.FTP.SavePath + ";" + "chmod 777 " + remoteFileName + ";" + "./" + remoteFileName
		log.Println(shellCmd)
		err = session.Run(shellCmd)
		if err != nil {
			log.Fatal(err)
		}
	}()

	// 监听CTRL+C信号
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	for {
		sig := <-ch
		log.Println("sig", sig)
		switch sig {
		case os.Interrupt, syscall.SIGTERM:
			// 很幸运发现，Linux进程地名字就是正在运行地编译文件名字
			session, err := sshClient.NewSession()
			if err != nil {
				log.Println(err)
				return err
			}
			Shell := "pkill -f helloweb"
			log.Println("Shell =", Shell)
			err = session.Run(Shell)
			if err != nil {
				log.Println(err)
				return err
			}
			defer session.Close()
		}
	}
}

// uploadFile 上传文件到服务器
func uploadFile(filePath string, sftpClient *sftp.Client) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer file.Close()
	remoteFileName := path.Base(filePath)
	log.Println(remoteFileName)
	dstFile, err := sftpClient.Create(path.Join(Conf.FTP.SavePath, remoteFileName))
	if err != nil {
		log.Println(err)
		return "", err
	}
	log.Println(dstFile)
	byteNum := 512000
	buf := make([]byte, byteNum)
	i := 0
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if n < byteNum {
			buff := make([]byte, 0)
			for a := 0; a < n; a++ {
				buff = append(buff, buf[a])
			}
			dstFile.Write(buff)
		} else {
			dstFile.Write(buf)
		}
		log.Println("文件包", i, "一次上传字节数", n)
		i++
	}
	defer dstFile.Close()
	return remoteFileName, nil
}
