package net_monitor

import (
	"common/config"
	"common/dbops"
	"common/taskrunner"
	"fmt"
	"github.com/shirou/gopsutil/net"
	"strconv"
	"sync"
	"time"
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
			if _, ok := ifaceMap[ioCounterStat.Name]; !ok {

				iface := &Iface{
					name:     ioCounterStat.Name,
					ip:       config.Conf.Interface[ioCounterStat.Name].IP,
					lastSend: ioCounterStat.BytesSent,
					lastRecv: ioCounterStat.BytesRecv,
					mtx:      &sync.Mutex{},
				}
				ifaceMap[ioCounterStat.Name] = iface
			} else {
				iface := ifaceMap[ioCounterStat.Name]

				//fmt.Printf(" iocounter.byterecv %d lastRecv %d iocounterstat %d lastSend %d \n", ioCounterStat.BytesRecv, iface.lastRecv, ioCounterStat.BytesSent, iface.lastSend)
				if ioCounterStat.BytesRecv < iface.lastRecv {
					iface.lastRecv = 0
				}
				if ioCounterStat.BytesSent < iface.lastSend {
					iface.lastSend = 0
				}
				received := ioCounterStat.BytesRecv - iface.lastRecv
				send := ioCounterStat.BytesSent - iface.lastSend

				iface.lastRecv = ioCounterStat.BytesRecv
				iface.lastSend = ioCounterStat.BytesSent

				t := time.Now()
				//ctime := t.Format("Jan 02 2006, 15:04:05")

				if received > 0 {
					recvRecord := new(dbops.MonitorRecord)
					recvRecord.Type = "net_recv"
					recvFloatValue := float64(received / 1024 / uint64(config.Conf.Global.Net_Update_Interval_Seconds))
					recvRecord.Value = strconv.FormatFloat(recvFloatValue, 'f', 6, 64)
					recvRecord.IP = iface.ip
					recvRecord.CreateTime = t

					if recvFloatValue > float64(config.Conf.Global.Net_Update_Threshold) {
						dc <- recvRecord
					}
				}

				if send > 0 {
					sendRecord := new(dbops.MonitorRecord)
					sendRecord.Type = "net_send"
					sendFloatValue := float64(send / 1024 / uint64(config.Conf.Global.Net_Update_Interval_Seconds))
					sendRecord.Value = strconv.FormatFloat(sendFloatValue, 'f', 6, 64)
					sendRecord.IP = iface.ip
					sendRecord.CreateTime = t

					if sendFloatValue > float64(config.Conf.Global.Net_Update_Threshold) {
						dc <- sendRecord
					}
				}
			}
		}
	}

	return nil
}

func Execute(record interface{}) error {
	r := record.(*dbops.MonitorRecord)
	displayTime := r.CreateTime.Format("2006-01-02, 15:04:05")

	fmt.Printf("Consume record: ip: %s, type: %s, value: %s, time: %s, properties: %s \n", r.IP, r.Type, r.Value, displayTime, r.Properties)
	err := dbops.InsertMonitorRecord(r)

	if err != nil {
		println("Net Execute record fail", err)
		return err
	}
	return nil
}
