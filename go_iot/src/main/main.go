package main

import (
	"common/config"
	"flag"
	"fmt"
	"log"
	"net_monitor"
	"os"
	"os/signal"
)

const (
	chanSize = 16
)

var (
	cfgFile = flag.String("config", `conf/sample.conf`, "Configuration file")
)

func init() {
	flag.Parse()
	if *cfgFile == "" {
		log.Fatal("Configuration file must be specified")
	}
}

func main() {
	fmt.Println("Go IOT monitor start...")

	cfg, err := config.NewConfig(*cfgFile)
	if err != nil {
		log.Fatal(err)
	}

	net_monitor.StartNewMonitor(cfg)

	sch := make(chan os.Signal)
	signal.Notify(sch, os.Interrupt, os.Kill)
	<-sch
}
