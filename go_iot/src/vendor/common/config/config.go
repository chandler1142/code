package config

import (
	"github.com/baidu/openedge/logger"
	cfg "gopkg.in/gcfg.v1"
)

var (
	defaultUpdateInterval uint = 600
)

type InterfaceDefinition struct {
	IP string
}

type Config struct {
	Global struct {
		Update_Interval_Seconds uint
	}
	Interface map[string]*InterfaceDefinition
}

func NewConfig(path string) (*Config, error) {
	var c Config
	c.Global.Update_Interval_Seconds = defaultUpdateInterval
	if err := cfg.ReadFileInto(&c, path); err != nil {
		logger.Warnf("Reading config path fail %v", err)
		return nil, err
	}

	return &c, nil
}
