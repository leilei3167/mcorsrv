package main

import (
	"database/sql"
	"errors"
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
	user := User{}
	// 1.获取第一条记录,first,take,last代表是否按照主键排序
	db.First(&user)
	// SELECT * FROM users ORDER BY id LIMIT 1;
	fmt.Println("#1 First:", user)
	// 获取一条记录，没有指定排序字段
	db.Take(&user, "age=?", 18)
	// SELECT * FROM users LIMIT 1;
	fmt.Println("#2 Take:", user)
	// 获取最后一条记录（主键降序）
	db.Last(&user, "age=?", 18)
	// SELECT * FROM users ORDER BY id DESC LIMIT 1;
	fmt.Println("#3 Last:", user)

	result := db.First(&user, "age=?", 18)
	// 检查 ErrRecordNotFound 错误
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		fmt.Println("没有找到数据")
	} else {
		fmt.Println("#4 First:", user)
	}

	//2.检索全部对象
	var usrs []User
	db.Find(&usrs)
	for _, usr := range usrs {
		fmt.Println("Find:", usr)

	}

	//3.条件查询,使用string来拼接语句,常用,会自动从user实例中识别表,如果不传入实例,则需要指定table
	db.Where("name=leilei").First(&user)        //name必须是表中的列名
	db.Where(&User{Name: "hdhsa"}).First(&user) //建议的形式,统一都用Go语言中的字段名,屏蔽细节
	//Find不会像take first last这样加limit,将会返回所有符合条件的数据

}
