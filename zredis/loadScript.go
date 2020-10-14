package zredis

import (
	"gitlab.icsoc.net/cc/gocommon/logger"
	"strings"
	"sync"
)

type ScriptLoader struct {
	// lock
	mu *sync.RWMutex
	// 注册成功之后保存key为script value为sha
	scriptSha map[string]string
	// 实例
	client Redis
}

// 注册脚本并加载
func NewLoader(redisName string, scripts []string) *ScriptLoader {

	var l = ScriptLoader{
		mu:        new(sync.RWMutex),
		scriptSha: make(map[string]string, 10),
		client:    GetRedisClient(redisName),
	}
	register(&l, scripts)
	return &l
}

// 获取sha
func (l *ScriptLoader) GetSha(key string) string {
	// 从map中获取
	l.mu.RLock()
	sha := l.scriptSha[key]
	l.mu.RUnlock()
	// 检查redis中是否存在此脚本和内存中是否存在
	if res, err := l.client.ScriptExists(key).Result(); err != nil || strings.EqualFold("", sha) || len(res) == 0 || !res[0] {
		//加载脚本
		sha = loadOneLuaScript(l.client, key)
		l.mu.Lock()
		l.scriptSha[key] = sha
		l.mu.Unlock()
	}
	return sha
}

// 加载脚本
func register(l *ScriptLoader, scripts []string) {
	retry := 3
	for i := 0; i < retry; i++ {
		// 判断是否存在
		bools, loadErr := l.client.ScriptExists(scripts...).Result()
		// 如果调用失败
		if loadErr != nil {
			logger.Errorf("LoadLuaScript err:%v", loadErr)
		} else {
			for i, b := range bools {
				if !b {
					l.mu.Lock()
					l.scriptSha[scripts[i]] = loadOneLuaScript(l.client, scripts[i])
					l.mu.Unlock()
				}
			}
			logger.Infof("loadscript success")
			break
		}
	}
}

func loadOneLuaScript(c Redis, s string) string {
	retry := 3
	res := ""
	for i := 0; i < retry; i++ {
		sha, err := c.ScriptLoad(s).Result()
		if err != nil {
			logger.Errorf("load %s error: %v retry: %d", s, err, i)
		} else {
			res = sha
			break
		}
	}
	return res
}
