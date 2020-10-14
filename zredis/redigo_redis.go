package zredis

import (
	"github.com/go-redis/redis/v7"
	"gitlab.icsoc.net/cc/gocommon/config"
)

type GoRedisClient struct {
	*redis.Client
	config *config.RedisConfig
}

func NewGoRedisClient(conf *config.RedisConfig) (Redis, error) {
	options := &redis.Options{
		Addr:         conf.Addr,
		Password:     conf.Password,
		DB:           conf.Db,
		MinIdleConns: conf.MinIdleConns,
	}
	// TODO: 是否提供redis clients 监听设置?
	client := redis.NewClient(options)
	if _, err := client.Ping().Result(); err != nil {
		return nil, err
	}

	return &GoRedisClient{
		client,
		conf,
	}, nil
}
