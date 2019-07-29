package config

import (
	"github.com/baidu/openedge/logger"
	cfg "gopkg.in/gcfg.v1"
)

var (
	defaultUpdateInterval uint = 600
	Conf                  Config
)

type InterfaceDefinition struct {
	IP string
}

type Config struct {
	Global struct {
		Net_Update_Interval_Seconds uint
	}
	Interface map[string]*InterfaceDefinition
}

func NewConfig(path string) {
	Conf.Global.Net_Update_Interval_Seconds = defaultUpdateInterval
	if err := cfg.ReadFileInto(&Conf, path); err != nil {
		logger.Warnf("Reading config path fail %v", err)
	}
}
