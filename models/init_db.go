package models

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func InitDB() {
	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("viper read config failed: %v", err)
	}
	host := viper.GetString("mysql.host")
	port := viper.GetString("mysql.port")
	user := viper.GetString("mysql.user")
	password := viper.GetString("mysql.password")
	dbname := viper.GetString("mysql.dbname")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, dbname)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt:            true, //缓存预编译命令
		SkipDefaultTransaction: true, //禁用默认事务操作
	})
	if err != nil {
		log.Fatalf("gorm open failed: %v", err)
	}
	DB = db

	err = DB.AutoMigrate(&UserInfo{}, &Video{}, &Comment{}, &UserLogin{})
	if err != nil {
		log.Fatalf("auto migrate failed: %v", err)
	}
}
