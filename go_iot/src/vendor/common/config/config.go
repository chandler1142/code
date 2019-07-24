package config

import (
	"fmt"
	"github.com/baidu/openedge/logger"
	cfg "gopkg.in/gcfg.v1"
)

func init() {
	fmt.Println("init common config")
}

var (
	defaultUpdateInterval uint = 600
	Conf                  Config
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

func NewConfig(path string) {
	Conf.Global.Update_Interval_Seconds = defaultUpdateInterval
	if err := cfg.ReadFileInto(&Conf, path); err != nil {
		logger.Warnf("Reading config path fail %v", err)
	}
}
