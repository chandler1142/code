package net_monitor

import (
	"common/taskrunner"
	"fmt"
	"github.com/shirou/gopsutil/net"
)

func init()  {
	fmt.Println("init net main")
}

func Dispatch(dc taskrunner.DataChan) error {

	ioCounterStats, err := net.IOCounters(true)
	if err != nil {
		fmt.Println(err)
	}

	for _, ioCounterStat := range ioCounterStats {
		fmt.Println(ioCounterStat)
	}


	//Get interface name and ip by this API
	//stats, err := net.Interfaces()
	//if err != nil {
	//	fmt.Printf("Get interface stats er: %v \n", err)
	//	return err
	//}
	//log.Println("Start to collect interfaces info...")
	//for _, stat := range stats {
	//	log.Println(stat)
	//}

	return nil
}

func Execute(dc taskrunner.DataChan) error {
	fmt.Println("Execute...")
	return nil
}
