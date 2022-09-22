package initialize

import (
	"encoding/json"
	"fmt"

	"mxshop_api/user_web/global"

	"github.com/spf13/viper"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"

	"go.uber.org/zap"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func InitConfig() {
	release := GetEnvInfo("MXSHOP_DEBUG")
	configFIlePrefix := "config"
	configFileName := fmt.Sprintf("%s-pro.yaml", configFIlePrefix)
	if !release {
		configFileName = fmt.Sprintf("%s-debug.yaml", configFIlePrefix)
	}
	v := viper.New()
	v.SetConfigFile(configFileName)

	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	// 这个对象应该设置为全局变量

	if err := v.Unmarshal(&global.NacosConfig); err != nil {
		panic(err)
	}
	zap.S().Debugf("nacos配置信息:%#v", global.NacosConfig)

	/*	v.WatchConfig()
		v.OnConfigChange(func(e fsnotify.Event) {
			zap.S().Info("配置文件产生变化: %v", e.Name)
			_ = v.ReadInConfig()
			_ = v.Unmarshal(global.ServerConfig)
			zap.S().Infof("变更后的配置信息: %v", global.ServerConfig)
		})*/

	// 从nacos中读取配置信息
	sc := []constant.ServerConfig{
		{
			IpAddr: global.NacosConfig.Host,
			Port:   global.NacosConfig.Port,
		},
	}

	cc := constant.ClientConfig{
		NamespaceId:         global.NacosConfig.Namespace, // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",

		LogLevel: "debug",
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		panic(err)
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group,
	})
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(content), &global.ServerConfig)
	if err != nil {
		zap.S().Fatalf("读取nacos配置失败： %s", err.Error())
	}

	zap.S().Debugf("配置信息:%#v", global.ServerConfig)
	// 想要将一个json字符串转换成struct，需要去设置这个struct的tag
	fmt.Printf("from nacos:%#v", content)
}
