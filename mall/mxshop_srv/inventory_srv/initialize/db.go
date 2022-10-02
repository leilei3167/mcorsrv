package initialize

import (
	"fmt"
	"log"
	"os"
	"time"

	"mxshop_srv/inventory_srv/global"
	"mxshop_srv/inventory_srv/model"

	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func InitDB() {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取dsn详情
	c := global.ServerConfig.MysqlInfo
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.Name)

	// 定义全局的sql,这样能够将某些慢sql打印出来,便于debug
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  true,        // 禁用彩色打印
		},
	)
	var err error
	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                 newLogger,
		NamingStrategy:         schema.NamingStrategy{SingularTable: true}, // 单数命名表
		SkipDefaultTransaction: false,
	})
	if err != nil {
		panic(err)
	}

	d, _ := global.DB.DB()
	d.SetMaxOpenConns(100)
	d.SetMaxIdleConns(100)
	d.SetConnMaxIdleTime(time.Hour)

	_ = global.DB.AutoMigrate(&model.Inventory{})

	client := goredislib.NewClient(&goredislib.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})
	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)
	global.RS = redsync.New(pool)
}
