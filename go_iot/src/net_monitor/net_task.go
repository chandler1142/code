package net_monitor

import (
	"common/config"
	"common/taskrunner"
	"fmt"
	"github.com/shirou/gopsutil/net"
	"sync"
)

var (
	ifaceMap = make(map[string]*Iface)
)

func Dispatch(dc taskrunner.DataChan) error {

	ioCounterStats, err := net.IOCounters(true)
	if err != nil {
		fmt.Println(err)
	}

	for _, ioCounterStat := range ioCounterStats {
		if config.Conf.Interface[ioCounterStat.Name] != nil {
			fmt.Println(ioCounterStat)
			if _,ok := ifaceMap[ioCounterStat.Name]; !ok {
				iface := &Iface{
					name: ioCounterStat.Name,
					ip: config.Conf.Interface[ioCounterStat.Name].IP,
					lastSend: ioCounterStat.BytesSent,
					lastRecv: ioCounterStat.BytesRecv,
					mtx: &sync.Mutex{},
				}
				ifaceMap[ioCounterStat.Name] = iface
			} else {
				iface := ifaceMap[ioCounterStat.Name]
				received := ioCounterStat.BytesRecv - iface.lastRecv
				send := ioCounterStat.BytesSent - iface.lastSend
				iface.lastRecv = ioCounterStat.BytesRecv
				iface.lastSend = ioCounterStat.BytesSent
				fmt.Printf("received: %d, send: %d \n", received, send)

				fmt.Printf("received rate: %.2f kb/s, send rate: %.2f kb/s \n ", float64(received/1024/5), float64(send/1024/5))
			}
		}

	}
	fmt.Println("==============================================")

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
