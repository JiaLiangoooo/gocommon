package zredis

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/pkg/errors"
	"gitlab.icsoc.net/cc/gocommon/config"
	"gitlab.icsoc.net/cc/gocommon/logger"
)

var (
	redisClients = make(map[string]Redis)
)

func init() {
	for name, conf := range config.GetMutiRedisConfig() {
		if err := SetRedisClient(name, conf); err != nil {
			logger.Errorf("set redis client (%s) conf(%#v) error: %v", name, conf, err)
		}
		logger.Infof("load redis: %s, %#v", name, conf)
	}
}

type Redis redis.Cmdable

//type Redis interface {
//	redis.Cmdable
//}

// GetRedisClient 获取当前的redis
func GetRedisClient(name string) Redis {
	if len(name) == 0 {
		logger.Errorf("获取redis失败, 传入name为空")
		return nil
	}
	r, ok := redisClients[name]
	if !ok {
		logger.Errorf("redis client (%s) does not exist,  please make sure you are regist it ", name)
	}
	return r
}

// SetRedisClient 设置RedisClients
func SetRedisClient(name string, c *config.RedisConfig) error {
	client, err := NewGoRedisClient(c)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("初始化 redis client %s error, %#v", name, c))
	}
	redisClients[name] = client
	return err
}
