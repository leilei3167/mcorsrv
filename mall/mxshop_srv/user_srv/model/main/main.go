package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
	"user_srv/user_srv/model"
	"user_srv/user_srv/pkg/password"
)

func main() {

	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取dsn详情
	dsn := "root:123456@tcp(127.0.0.1:3306)/mxshop_user_srv?charset=utf8mb4&parseTime=True&loc=Local"

	//定义全局的sql,这样能够将某些慢sql打印出来,便于debug
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  true,        // 禁用彩色打印
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:         newLogger,
		NamingStrategy: schema.NamingStrategy{SingularTable: true}, //单数命名表
	})
	if err != nil {
		panic(err)
	}
	_ = db.AutoMigrate(&model.User{})

	pwd := "admin123"

	for i := 0; i < 10; i++ {
		p, _ := password.Encode(pwd)
		user := model.User{
			NickName: fmt.Sprintf("leilei_%d", i),
			Mobile:   fmt.Sprintf("186028431%d", i),
			Password: p,
		}
		db.Save(&user)
	}

}
