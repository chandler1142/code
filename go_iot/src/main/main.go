package main

import (
	"common/config"
	"disk_monitor"
	"flag"
	"fmt"
	"log"
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

	config.NewConfig(*cfgFile)

	//net_monitor.StartNewMonitor()
	disk_monitor.StartDiskMonitor()

	sch := make(chan os.Signal)
	signal.Notify(sch, os.Interrupt, os.Kill)
	<-sch
}


