/*
@time : 2020/7/14 9:17 上午
@author : zcq
@description :鉴权测试
*/
package auth

import (
	"github.com/stretchr/testify/assert"
	"gitlab.icsoc.net/cc/gocommon/config"
	"gitlab.icsoc.net/cc/gocommon/zredis"
	"log"
	"testing"
)

var (
	NAME         = "auth"
	REDIS_CONFIG = &config.RedisConfig{
		Addr:         "127.0.0.1:6379",
		Password:     "",
		Db:           0,
		MinIdleConns: 10,
	}
)

//初始化配置
func initConfig() *config.ServiceConfig {
	authConfig := &config.AuthConfig{Host: "https://account-test.icsoc.net"}
	serviceConfig := config.ServiceConfig{AuthConfig: authConfig}
	config.SetServiceConfig(serviceConfig)
	zredis.SetRedisClient(NAME, REDIS_CONFIG)
	return &serviceConfig
}

//测试正常的token
func TestNormalAuth(t *testing.T) {
	token := "4a39108cabfbe55f23e2a33529f55cbb563dda00"
	serviceConfig := initConfig()
	log.Printf("验证token:%s,host:%s", token, serviceConfig.AuthConfig.Host)
	user, err := CheckAuth(token)
	assert.Nil(t, err)
	log.Println(user)
}

//测试没有配置的情况
func TestNoConfigAuth(t *testing.T) {
	token := "30a2805080b8d626a8800a9268b1936ae281acc1"
	user, err := CheckAuth(token)
	assert.Nil(t, err)
	if user != nil {
		t.Fail()
	}
}
