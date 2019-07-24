package net_monitor

import (
	"common/config"
	"common/taskrunner"
	"fmt"
)

func StartNewMonitor() {
	//1. read conf file and init the monitor map
	if len(config.Conf.Interface) == 0 {
		fmt.Println("No interfaces specified")
		return
	}

	//2. init the task runner and start the goroutine worker
	r := taskrunner.NewRunner("NetInterfaceMonitor", 3, true, Dispatch, Execute)
	w := taskrunner.NewWorker(5, r)
	go w.StartWorker()

}
