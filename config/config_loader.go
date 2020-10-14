package config

import (
	"github.com/JiaLiangoooo/gocommon/logger"
	"os"
)

var (
	serviceConfig *ServiceConfig
)

func init() {
	var configFile string

	configFile = os.Getenv(CONF_SERVICE_FILE_PATH)
	if err := ServiceConfigInit(configFile); err != nil {
		// %#v 输出 Go 语言语法格式的值
		logger.Errorf("[serviceConfigInit] %#v", err)
		serviceConfig = nil
	}
}

func Load() {
	loadServiceConfig()
}

// GetServiceConfig 获取service配置
func GetServiceConfig() *ServiceConfig {
	return serviceConfig
}

func GetHttpConfig() *HttpConfig {
	if serviceConfig == nil {
		return nil
	}
	return serviceConfig.HttpConfig
}

func GetMutiRedisConfig() map[string]*RedisConfig {
	if serviceConfig == nil {
		return nil
	}
	return serviceConfig.MutiRedis
}

func GetRedisConfig(name string) *RedisConfig {
	if serviceConfig == nil {
		return nil
	}
	return serviceConfig.MutiRedis[name]
}

func GetRabbitMqConfig() *RabbitMqConfig {
	if serviceConfig == nil {
		return nil
	}
	return serviceConfig.RabbitMqConfig
}

//获取token认证配置
func GetAuthConfig() *AuthConfig {
	if serviceConfig == nil {
		return nil
	}
	return serviceConfig.AuthConfig
}
