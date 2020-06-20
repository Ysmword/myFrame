package common

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// 远程开发调试
// 生成最新编译文件
// 将编译文件上传到远程服务器
// 在服务器上运行

// Connect 创建sftp连接
func Connect() (*ssh.Session, *sftp.Client, error) {

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

	session, err := sshClient.NewSession()
	if err != nil {
		log.Println(err)
		return nil,nil, err
	}

	return session, sftpClient, nil
}

// UploadRun 上传可运行文件和运行文件
func UploadRun() error {

	session,sftpClient, err := Connect()
	if err != nil {
		log.Println(err)
		return err
	}
	defer sftpClient.Close()
	defer session.Close()

	// 先生成可执行文件
	cmdShell := "cd " + HomePath + ";set CGO_ENABLED=0;set GOOS=linux;set GOARCH=amd64;go build"
	cmd := exec.Command("cmd", "/C", cmdShell)
	data, err := cmd.Output()
	if err != nil {
		log.Println(err)
		return err
	}

	if len(data) != 0 {
		log.Println(string(data))
		return errors.New(string(data))
	}

	// 上传可执行文件
	localFilePath := HomePath + "/helloweb"
	if !IsExist(localFilePath) {
		log.Println(HomePath + "/helloweb" + "文件不存在")
		return errors.New(HomePath + "/helloweb" + "文件不存在")
	}
	file, err := os.Open(localFilePath)
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()

	remoteFileName := path.Base(localFilePath)
	log.Println(remoteFileName)
	dstFile, err := sftpClient.Create(path.Join(Conf.FTP.SavePath, remoteFileName))
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(dstFile)

	defer dstFile.Close()

	buf := make([]byte, 1024)
	for {
		n, _ := file.Read(buf)
		if n == 0 {
			break
		}
		dstFile.Write(buf)
	}

	// 运行命令
	session.Stderr = os.Stderr
	session.Stdout = os.Stdout
	shellCmd := "cd "+ Conf.FTP.SavePath+";"+"chmod 777 "+remoteFileName+";"+"./"+remoteFileName
	log.Println(shellCmd)
	err = session.Run(shellCmd)
	if err!=nil{
		log.Println(err)
		return err
	}
	return nil
}
