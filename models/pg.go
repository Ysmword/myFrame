package models

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"log"

	"helloweb/common"
)

// db postgres数据库的变量
var db *sql.DB

// InitDB 初始化数据库
func InitDB() {
	if common.Conf == nil {
		log.Fatal("conf.Conf is nil")
	}
	if common.Conf.Postgres.Dbname == "" || common.Conf.Postgres.Host == "" || common.Conf.Postgres.Password == "" || common.Conf.Postgres.Sslmode == "" || common.Conf.Postgres.User == "" {
		log.Fatal("illegal database connect info")
	}
	if common.Conf.Postgres.Port < 0 {
		log.Fatal("Conf.Postgres.Port is illegal")
	}

	// 打开数据库的信息
	dbInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		common.Conf.Postgres.Host,
		common.Conf.Postgres.Port,
		common.Conf.Postgres.User,
		common.Conf.Postgres.Password,
		common.Conf.Postgres.Dbname,
		common.Conf.Postgres.Sslmode,
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
