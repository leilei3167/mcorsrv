package initialize

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger() {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, _ := config.Build()
	zap.ReplaceGlobals(logger) //替换全局的logger
	//S()可以获得一个全局的sugar,L()则是全局的Logger
	//S()和L()提供一个安全的全局日志的使用
}
