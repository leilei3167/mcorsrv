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

type User struct {
	ID           uint
	Name         string
	Email        *string //设置为可以设置为零值
	Age          uint8
	Birthday     *time.Time
	MemberNumber sql.NullString //设置为可以设置为零值
	ActivatedAt  sql.NullTime
	CreatedAt    time.Time
	UpdatedAt    time.Time
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
	_ = db.AutoMigrate(&User{})
	//1.单条插入
	user := User{Name: "Leilei"}
	result := db.Create(&user)
	fmt.Printf("id: %v err:%v rows affected:%v\n", user.ID, result.Error, result.RowsAffected)
	//2.批量插入(gorm会将大量数据生成单一的sql语句,性能比较高)
	var users = []User{{Name: "jinzhu1"}, {Name: "jinzhu2"}, {Name: "jinzhu3"}}
	//	db.Create(&users)

	//3.CreateInBatches 分批次插入,当数量超大时,由于sql有长度限制,上述方法就无法使用,需要分批次插入
	db.CreateInBatches(users, 2) //size代表每次最多提交多少个

	//4.map创建,更灵活,不用事先创建结构体
	db.Model(&User{}).Create(map[string]interface{}{"Name": "from map", "Age": 18})
	//5.批量map创建
	db.Model(&User{}).Create([]map[string]interface{}{
		{"Name": "jinzhu_1", "Age": 18},
		{"Name": "jinzhu_2", "Age": 20},
	})

}
