package config

import (
	"github.com/creasty/defaults"
)

//认证配置
type AuthConfig struct {
	Host string `yaml:"host" json:"host"`
}

func (c *AuthConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if err := defaults.Set(c); err != nil {
		return err
	}
	type plain AuthConfig
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}
	return nil
}
