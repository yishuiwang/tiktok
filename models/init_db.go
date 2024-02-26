package models

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"tiktok/config"
	"tiktok/logger"
)

var DB *gorm.DB
var Redis *redis.Client

func InitRedis() {
	ctx := context.Background()
	host := config.GetConfig("redis.host")
	port := config.GetConfig("redis.port")
	password := config.GetConfig("redis.password")
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password, // no password set
		DB:       0,        // use default DB
	})
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		logger.ZapLogger.Error("connect redis failed", logger.Error(err))
	}
	logger.ZapLogger.Info("connect redis success", logger.String("host", host), logger.String("port", port))
	Redis = rdb

}

func InitMysql() {
	host := config.GetConfig("mysql.host")
	port := config.GetConfig("mysql.port")
	user := config.GetConfig("mysql.user")
	password := config.GetConfig("mysql.password")
	dbname := config.GetConfig("mysql.dbname")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, dbname)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt:            true, //缓存预编译命令
		SkipDefaultTransaction: true, //禁用默认事务操作
	})
	if err != nil {
		logger.ZapLogger.Error("connect mysql failed", logger.Error(err))
	}
	logger.ZapLogger.Info("connect mysql success")
	DB = db
	//
	//err = DB.AutoMigrate(&UserInfo{}, &Video{}, &Comment{}, &UserLogin{})
	//if err != nil {
	//	log.Fatalf("auto migrate failed: %v", err)
	//}
}
