package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"smServer/config"
	"xorm.io/xorm"
)

var Engin *xorm.Engine

func TestDb() {
	mysqlConfig, err := config.File.GetSection("mysql")
	if err != nil {
		log.Println(err)
		panic(err)
	}
	dbConn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		mysqlConfig["user"],
		mysqlConfig["password"],
		mysqlConfig["host"],
		mysqlConfig["port"],
		mysqlConfig["dbname"])

	fmt.Println(dbConn)

	// var err error
	Engin, err = xorm.NewEngine("mysql", dbConn)
	if err != nil {
		log.Println("数据库连接失败：", err)
		panic(err)
	}
	err = Engin.Ping()
	if err != nil {
		log.Println("数据库ping不通", err)
		panic(err)
	}
	maxIdel := config.File.MustInt("mysql", "max_idel", 2)
	maxConn := config.File.MustInt("mysql", "max_conn", 1)
	Engin.SetMaxIdleConns(maxIdel)
	Engin.SetMaxOpenConns(maxConn)
	Engin.ShowSQL()
	log.Println("数据化初始化完成")
}
