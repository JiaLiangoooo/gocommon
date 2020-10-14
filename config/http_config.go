package config

import (
	"github.com/creasty/defaults"
	"time"
)

type HttpConfig struct {
	Host string        `yaml:"host" json:"host"`
	Port int           `yaml:"port" json:"port"`
	TTL  time.Duration `default:"10s" yaml:"ttl" json:"ttl"`
}

func (c *HttpConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if err := defaults.Set(c); err != nil {
		return err
	}
	type plain HttpConfig
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}
	return nil
}
