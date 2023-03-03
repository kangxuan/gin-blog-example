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
	CreatedOn  int `gorm:"autoCreateTime" json:"created_on"`
	ModifiedOn int `gorm:"autoUpdateTime" json:"modified_on"`
}

func SetUp() {
	var err error
	// 如果数据库类型为MySQL
	if settings.DatabaseSetting.Type == "mysql" {
		db, err = gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			settings.DatabaseSetting.User,
			settings.DatabaseSetting.Password,
			settings.DatabaseSetting.Host,
			settings.DatabaseSetting.Port,
			settings.DatabaseSetting.Name),
		), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   settings.DatabaseSetting.TablePrefix, // 表前缀
				SingularTable: true,                                 // 禁用复数表名
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
