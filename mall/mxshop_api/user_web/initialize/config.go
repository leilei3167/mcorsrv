package initialize

import (
	"fmt"
	"github.com/spf13/viper"
	"user_api/user_web/global"

	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func InitConfig() {
	debug := GetEnvInfo("MXSHOP_DEBUG")
	configFIlePrefix := "config"
	configFileName := fmt.Sprintf("%s-pro.yaml", configFIlePrefix)
	if debug {
		configFileName = fmt.Sprintf("%s-debug.yaml", configFIlePrefix)
	}
	v := viper.New()
	v.SetConfigFile(configFileName)

	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	//这个对象应该设置为全局变量

	if err := v.Unmarshal(&global.ServerConfig); err != nil {
		panic(err)
	}
	zap.S().Infof("配置信息:%v", global.ServerConfig)

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		zap.S().Info("配置文件产生变化: %v", e.Name)
		_ = v.ReadInConfig()
		_ = v.Unmarshal(global.ServerConfig)
		zap.S().Infof("变更后的配置信息: %v", global.ServerConfig)
	})

}
