package net_monitor

import (
	"common/config"
	"common/taskrunner"
	"fmt"
	"time"
)

func StartNewMonitor() {
	//1. read conf file and init the monitor map
	if len(config.Conf.Interface) == 0 {
		fmt.Println("No interfaces specified")
		return
	}

	//2. init the task runner and start the goroutine worker
	r := taskrunner.NewRunner("NetInterfaceMonitor", 64, true, Dispatch, Execute)
	w := taskrunner.NewWorker(time.Duration(config.Conf.Global.Net_Update_Interval_Seconds), r)
	go w.StartDispatch()
	go w.StartExecutor()
}
