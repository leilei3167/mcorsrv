package main

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type Product struct {
	gorm.Model //继承gorm的默认公共字段,已经自带了ID主键,创建修改时间
	Code       sql.NullString
	Price      uint
}

func main() {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取dsn详情
	dsn := "root:123456@tcp(127.0.0.1:3306)/gorm_test?charset=utf8mb4&parseTime=True&loc=Local"

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
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

	//根据实例生成表结构
	_ = db.AutoMigrate(&Product{}) //此处会有sql语句

	//CRUD
	//C
	db.Create(&Product{Code: sql.NullString{String: "leilei", Valid: true}, Price: 10000})

	//R
	var p Product                    //容纳结果
	db.First(&p, 1)                  //默认根据主键查找
	db.First(&p, "code=?", "leilei") //根据自定义的条件
	fmt.Println("查询p:", p)

	//U
	db.Model(&p).Update("Price", 200)                                                        //注意是结构体的字段名                               //更新单个字段,指名要更新的字段以及值
	db.Model(&p).Updates(Product{Price: 200, Code: sql.NullString{String: "", Valid: true}}) //Code 设置为0值则不会被更新              //仅更新非0字段
	//db.Model(&p).Updates(map[string]interface{}{"Price": 300, "Code": ""}) //以map的形式更新

	//D
	db.Delete(&p, 1)
}
