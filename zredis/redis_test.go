package zredis

import (
	"testing"
)
import (
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/JiaLiangoooo/gocommon/config"
)

var (
	name        = "test"
	redisConfig = &config.RedisConfig{
		Addr:         "127.0.0.1:6379",
		Password:     "",
		Db:           0,
		MinIdleConns: 10,
	}
)

func TestSetRedisClient(t *testing.T) {

	err := SetRedisClient(name, redisConfig)
	assert.Nil(t, err)

	client := GetRedisClient(name)
	assert.NotNil(t, client)

	key, value := "test", "1111"
	client.Set(key, value, 0)
	s := client.Get("test").Val()
	assert.Equal(t, value, s)
}
