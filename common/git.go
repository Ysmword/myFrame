package common


import (
	"fmt"
	"os/exec"
	"runtime"
	"log"
	"bufio"
	"html/template"
	"io/ioutil"
	"os"
	"strings"
)
// 这里实现对git的操作基础操作   这样封装，就不需要每个地方都要调用go-git包，我们同意调用就行了

/* 
1、需求就是能够：自己写完代码之后，能够实现自动提交，自动提交，自动拉取
*/

var (
	// GOOS 当前运行的操作系统
	GOOS = runtime.GOOS 
	// GitShellFileName 存放git add .,git commit,git push的shell文件,window10条件下，要保证有git.exe文件，才能运行脚本文件
	GitShellFileName = "gitShell.sh"
)

// GitShellInfo gitShellFile模板信息
type GitShellInfo struct {
	// git shell脚本文件路径
	GitShellFilePath string `json:"gitShellFilePath"`
	// CommitInfo 提交信息
	CommitInfo string `json:"commitInfo"`
	// GitEmail 提交的用户邮箱
	GitEmail string `json:"gitEmail"`
	// GitName 提交用户名
	GitName string  `json:"gitName"`
}

// ReadGitShellTmp 读取gitShell.tmp文件
func ReadGitShellTmp() (string, error) {

	if HomePath == "" {
		err := fmt.Errorf("ReadGitShellTmp HomePath is null")
		log.Println(err)
		return "",err
	}

	tmpFileName:= "gitShell.tmp"
	tmpPath := HomePath + "/" + tmpFileName

	if !IsExist(tmpPath){
		err := fmt.Errorf("no file: "+tmpPath)
		log.Println(err)
		return "",err
	}

	shellTmp, err := ioutil.ReadFile(tmpPath)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return string(shellTmp), nil
}


// WriteGitShellFile 生成gitShellFile.sh temp 模板 返回生成文件之后文件路径
func WriteGitShellFile(temp string) (string,error) {

	if temp == ""{
		err := fmt.Errorf("WriteGitShellFile temp is null")
		log.Println(err)
		return "",err
	}

	if HomePath == "" {
		err := fmt.Errorf("ReadGitShellTmp HomePath is null")
		log.Println(err)
		return "",err
	}
	var err error
	targetFileName := "gitShell.sh"
	// 生成文件的目标文件
	targetFilePath := HomePath + "/" + targetFileName

	// 创建，覆盖，可读可写
	f, err := os.OpenFile(targetFilePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
	if err != nil {
		log.Println(err)
		return "",err
	}

	// 这里不关闭，就会不能执行里面的文件
	defer f.Close()

	var commitInfo string

	fmt.Println("请输入提交信息")
	inputReader := bufio.NewReader(os.Stdin)
	// ReadString 以换行符号结束,输出会包含换行符号
	commitInfo, err = inputReader.ReadString('\n')
	if err != nil {
		log.Println(err)
		return "",err
	}
	commitInfo = strings.ReplaceAll(commitInfo, "\n", "")
	log.Println(commitInfo)

	// 创建一个模板对象
	gitTemp := template.New("gitTemp")

	// 解析模板
	gitTemp, err = gitTemp.Parse(temp)
	if err != nil {
		log.Println(err)
		return "",err
	}


	if Conf.Git.GitEmail == "" {
		err := fmt.Errorf("Conf.Git.GitEmail is null")
		log.Println(err)
		return "",err
	}

	if Conf.Git.GitName == "" {
		err := fmt.Errorf("Conf.Git.GitName is null")
		log.Println(err)
		return "",err
	}

	

	// 实例化一个GitShellInfo对象  更本地项目连接在一起
	gitShellInfo := &GitShellInfo{
		GitShellFilePath: HomePath,
		 CommitInfo: commitInfo,
		 GitEmail: Conf.Git.GitEmail,
		 GitName: Conf.Git.GitName,
		}

	// 模板填充
	err = gitTemp.Execute(f, gitShellInfo)
	if err != nil {
		log.Println(err)
		return "",err
	}
	return targetFilePath,nil
}


// GitCmd 执行gitShellFile文件
func GitCmd()error{
	
	log.Println("当前的操作系统",runtime.GOOS)

	if runtime.GOOS!="linux"{
		err := fmt.Errorf("暂时不支持linux系统意外的操作系统")
		log.Println(err)
		return err
	}

	// 读取gitShellTmp 文件
	temp, err := ReadGitShellTmp()
	if err != nil {
		log.Println(err)
		return err
	}

	// 根据模板写入文件
	gitShellFilePath,err := WriteGitShellFile(temp)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("gitShellFilePath:",gitShellFilePath)
	// 执行文件
	cmd := exec.Command("/bin/bash","-c",gitShellFilePath)
	log.Println(cmd.Args)

	data, err := cmd.Output()
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(string(data))
	return nil
}




