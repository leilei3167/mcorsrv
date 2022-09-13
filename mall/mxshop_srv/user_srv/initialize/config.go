package initialize

import (
	"fmt"
	"user_srv/user_srv/global"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
	//刚才设置的环境变量 想要生效 我们必须得重启goland
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
