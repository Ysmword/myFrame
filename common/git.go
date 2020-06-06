package common


import (
	"fmt"
	"os/exec"
	"runtime"

	"log"

	"github.com/go-git/go-git/v5"
)


// 这里实现对git的操作基础操作   这样封装，就不需要每个地方都要调用go-git包，我们同意调用就行了

/* 
1、需求就是能够：自己写完代码之后，能够实现自动提交，自动提交，自动拉取

*/

/*
git clone
git commit
git pull
git push
git branch
*/

// GitClone 克隆操作 lcoalPath 拷贝到具体地址目录 isBare 是否是裸项目 cloneOptions 
func GitClone(localPath string,isBare bool,url string,cloneOptions *git.CloneOptions)error{

	if localPath == ""{
		err := fmt.Errorf("GitClone localPath is null")
		log.Println(err)
		return err
	}

	isDir,err := IsDir(localPath)
	if err!=nil{
		log.Println(err)
		return err
	}
	
	if !isDir{
		err := fmt.Errorf("no floder in project " + localPath)
		log.Println(err)
		return err
	}

	if url == ""{
		err := fmt.Errorf("GitClone url is null")
		log.Println(err)
		return err
	}
	log.Println("拷贝仓库")
	_,err = git.PlainClone(localPath,isBare,cloneOptions)

	if err!=nil{
		log.Println(err)
		return err
	}
	return nil
}

var (
	// GOOS 当前运行的操作系统
	GOOS = runtime.GOOS 
	// GitShellFileName 存放git add .,git commit,git push的shell文件,window10条件下，要保证有git.exe文件，才能运行脚本文件
	GitShellFileName = "gitShell.sh"
)



// GitCmd 执行gitShellFile文件
func GitCmd()error{

	if HomePath == ""{
		err := fmt.Errorf("GitCmd HomePath is null")
		log.Println(err)
		return err
	}
	var gitShellFilePath = HomePath+"/"+GitShellFileName
	log.Println("gitShellFilePath:",gitShellFilePath)

	// 判断gitShellFilePath是文件是否存在
	if !IsExist(gitShellFilePath){
		err := fmt.Errorf("there is no such path: "+ gitShellFilePath)
		log.Println(err)
		return err
	}

	switch GOOS {
	case "windows":
		break
	case "linux":
		c := exec.Command("/bin/sh","-c",gitShellFilePath)
		if data,err := c.Output();err!=nil{
			log.Println(err)
			return err
		}else{
			log.Println(string(data))
		}
		break
	default:
		err := fmt.Errorf("只支持windows和Linux的操作系统")
		log.Println(err)
		return err
	}

	return nil
}




