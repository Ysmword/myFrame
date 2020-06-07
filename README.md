# myFrame


## 文件夹以及文件介绍
- common文件夹放置公共处理函数
- controller文件夹放置接口处理函数
- env文件夹放置启动项目函数
- models文件夹放置与数据库表对应的结构体
- tmpl 文件夹放置数据库表映射模板
- build.sh 部署文件
- conf.ini 配置文件
- gitShell.tmp 自动提交到git的模板文件
- gitShell.sh 有gitShell.tmp模板文件生成的自动提交到git的脚本文件
- modelGen.sh 将数据库表生成对应结构体的脚本文件

## 框架拥有的功能
- 配置文件读取
- 自动提交到git
- 发送短信
- 发送邮箱
- 静态服务器搭建（运行之后，可以直接访问：http://localhost:端口/file，就可以看到整个项目的静态文件）
- 一键生成数据库表的结构体映射文件
- 日志处理
