package models

import (
	"fmt"
	"gin-blog-example/settings"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
)

var db *gorm.DB

// Model 定义基础的模型字段
type Model struct {
	ID         int `gorm:"primary_key" json:"id"`
	CreatedOn  int `json:"created_on"`
	ModifiedOn int `json:"modified_on"`
}

// init 初始化数据库连接
func init() {
	var (
		err                                                     error
		dbType, dbName, user, password, host, port, tablePrefix string
	)

	sec, err := settings.Cfg.GetSection("database")
	if err != nil {
		log.Fatal(2, "Fail to get section 'database': %v", err)
	}

	dbType = sec.Key("TYPE").String()
	dbName = sec.Key("NAME").String()
	user = sec.Key("USER").String()
	password = sec.Key("PASSWORD").String()
	host = sec.Key("HOST").String()
	port = sec.Key("PORT").String()
	tablePrefix = sec.Key("TABLE_PREFIX").String()

	if dbType == "mysql" {
		db, err = gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, dbName)), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   tablePrefix, // 表前缀
				SingularTable: true,        // 禁用复数表名
			},
		})
	} else {
		log.Fatalln("db-type暂不支持")
		return
	}
	if err != nil {
		log.Fatalln(err)
	}
}
