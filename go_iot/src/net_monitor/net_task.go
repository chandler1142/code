package net_monitor

import (
	"common/taskrunner"
	"fmt"
	"github.com/shirou/gopsutil/net"
	"log"
)

func Dispatch(dc taskrunner.DataChan) error {
	//Get interface name and ip by this API
	stats, err := net.Interfaces()
	if err != nil {
		fmt.Printf("Get interface stats er: %v \n", err)
		return err
	}

	log.Println("Start to collect interfaces info...")
	for _, stat := range stats {
		log.Println(stat)
	}
	return nil
}

func Execute(dc taskrunner.DataChan) error {
	fmt.Println("Execute...")
	return nil
}
