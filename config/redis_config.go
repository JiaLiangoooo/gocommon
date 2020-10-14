package config

import "github.com/creasty/defaults"

type RedisConfig struct {
	Addr         string `required:"true" yaml:"addr" json:"addr"`
	Password     string `required:"true" yaml:"password" json:"password"`
	Db           int    `default:"0" yaml:"db" json:"db"`
	MinIdleConns int    `default:"20" yaml:"min_idle" json:"min_idle"`
}

func (c *RedisConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if err := defaults.Set(c); err != nil {
		return err
	}
	type plain RedisConfig
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}
	return nil
}
