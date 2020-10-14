package config

import "github.com/creasty/defaults"

type RabbitMqConfig struct {
	Uri        string `required:"true" yaml:"uri" json:"uri"`
	Exchange   string `required:"true" yaml:"exchange" json:"exchange"`
	Queue      string `required:"true" yaml:"queue" json:"queue"`
	Topic      string `required:"true" yaml:"topic" json:"topic"`
	RoutineKey string `yaml:"routine_key" json:"routine_key"`
}

func (c *RabbitMqConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if err := defaults.Set(c); err != nil {
		return err
	}
	type plain RabbitMqConfig
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}
	return nil
}
