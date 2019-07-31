package disk_monitor

import (
	"common/config"
	"common/dbops"
	"common/taskrunner"
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/disk"
	"strconv"
	"time"
)


func Dispatch(dc taskrunner.DataChan) error {

	//partitionStats, err := disk.Partitions(false)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(partitionStats)

	usageStat, err := disk.Usage("/")

	if err != nil {
		fmt.Println(err)
		return err
	}

	//fmt.Printf("total: %f, free: %f, used: %f usePersent: %f", float64(usageStat.Total/1024/1024/1024),float64(usageStat.Free/1024/1024/1024), float64(usageStat.Used/1024/1024/1024), usageStat.UsedPercent)

	properties, err := json.Marshal(usageStat)
	if err != nil {
		fmt.Println(err)
		return err
	}
	record := &dbops.MonitorRecord{
		Type: "disk",
		Value: strconv.FormatFloat(usageStat.UsedPercent, 'f', 6, 64),
		IP: config.Conf.Global.IP,
		CreateTime: time.Now(),
		Properties: string(properties),
	}

	dc <- record

	return nil
}

func Execute(record interface{}) error {
	fmt.Println(record)
	return nil
}

