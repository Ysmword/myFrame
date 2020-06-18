package common

import (
	"log"
	"os"
)

// 热编译
// 1、监听项目变化
// 2、关闭当前的进程，并且将当前的进程资源全部收回
// 3、重新编译文件
// 4、重新运行

// ProjectState 项目状态
type ProjectState struct{
	// NewDoc 新增文件
	NewDoc bool

	// ModFile 修改过文件
	ModFile bool

	// DelFile 删除文件
	DelFile bool
}

// IsHasMod 是否被修改过
func(ps *ProjectState)IsHasMod()bool{
	return ps.DelFile||ps.ModFile||ps.NewDoc
}


// Monitor 监听者
func Monitor(){
	log.Println(HomePath)
	// 查看整个项目文件是否有增加
	fileInfo,err :=os.Stat(HomePath)
	if err!=nil{
		log.Println(err)
		return
	}
	modTime := fileInfo.ModTime().String()

	// 已经修改了，则重新编译
	if modTime == LastModifiedTime {
		
	}
}