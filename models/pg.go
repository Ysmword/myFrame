package models

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"log"

	"helloweb/conf"
)

// db postgres数据库的变量
var db *sql.DB

// InitDB 初始化数据库
func InitDB()  {
	if conf.Conf == nil {
		log.Fatal("conf.Conf is nil")
	}
	if conf.Conf.Postgres.Dbname == "" || conf.Conf.Postgres.Host == "" || conf.Conf.Postgres.Password == "" || conf.Conf.Postgres.Sslmode == "" || conf.Conf.Postgres.User == "" {
		log.Fatal("illegal database connect info")
	}
	if conf.Conf.Postgres.Port < 0 {
		log.Fatal("Conf.Postgres.Port is illegal")
	}

	// 打开数据库的信息
	dbInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		conf.Conf.Postgres.Host,
		conf.Conf.Postgres.Port,
		conf.Conf.Postgres.User,
		conf.Conf.Postgres.Password,
		conf.Conf.Postgres.Dbname,
		conf.Conf.Postgres.Sslmode,
	)

	var err error
	//  打开数据库
	db, err = sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal("sql.Open error ", err)
	}
	// 查看是否连接成功
	if err := db.Ping(); err != nil {
		// 当连接不成功时，就要及时的关闭
		db.Close()
		err = fmt.Errorf("db.Ping error:%v", err)
		log.Fatal(err)
	}
	//  设置空闲连接池中最大的连接数
	db.SetMaxIdleConns(20)
	// 设置数据库的最大开放连接数，如果n≤0，则开放数量没有限制
	db.SetMaxOpenConns(0)
}
