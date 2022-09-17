package global

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"mxshop_srv/goods_srv/config"
)

var (
	// DB 全局的数据库实例,handler层直接依赖,应该考虑解耦,方便后期更换数据库
	DB           *gorm.DB
	ServerConfig *config.ServerConfig
	NacosConfig  = &config.NacosConfig{}
)

func init() {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取dsn详情
	dsn := "root:123456@tcp(127.0.0.1:3306)/mxshop_goods_srv?charset=utf8mb4&parseTime=True&loc=Local"

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
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:         newLogger,
		NamingStrategy: schema.NamingStrategy{SingularTable: true}, //单数命名表
	})
	if err != nil {
		panic(err)
	}

}
