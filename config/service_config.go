package config

import (
	"github.com/pkg/errors"
)

import (
	"github.com/JiaLiangoooo/gocommon/logger"
	"github.com/JiaLiangoooo/gocommon/yaml"
)

type ServiceConfig struct {
	HttpConfig     *HttpConfig             `yaml:"http" json:"http, omitempty"`
	RabbitMqConfig *RabbitMqConfig         `yaml:"rabbitmq" json:"rabbitmq, omitempty"`
	MutiRedis      map[string]*RedisConfig `yaml:"redis" json:"redis, omitempty"`
	AuthConfig     *AuthConfig             `yaml:"auth" json:"auth, omitempty"`
}

// ServiceConfigInit 初始化服务配置
func ServiceConfigInit(configFile string) error {
	if len(configFile) == 0 {
		return errors.Errorf("service config file name is nil")
	}
	serviceConfig = &ServiceConfig{}
	if _, err := yaml.UnmarshalYMLConfig(configFile, serviceConfig); err != nil {
		return err
	}
	logger.Debugf("加载service config配置文件,HttpConfig(%+v),RabbitMqConfig(%+v),MutiRedis(%v),AuthConfig(%+v)", serviceConfig.HttpConfig, serviceConfig.RabbitMqConfig, serviceConfig.MutiRedis, serviceConfig.AuthConfig)
	return nil
}

func loadServiceConfig() {

}

// SetServiceConfig 设置Service配置, 在设置完之后, 需要Load
func SetServiceConfig(config ServiceConfig) {
	serviceConfig = &config
}
