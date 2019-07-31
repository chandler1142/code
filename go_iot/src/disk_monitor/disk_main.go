package disk_monitor

import (
	"common/config"
	"common/taskrunner"
	"time"
)

func StartDiskMonitor() {
	//1. read conf file and init the monitor map


	//2. init the task runner and start the goroutine worker
	r := taskrunner.NewRunner("DiskInterfaceMonitor", 64, true, Dispatch, Execute)
	w := taskrunner.NewWorker(time.Duration(config.Conf.Global.Disk_Update_Interval_Seconds), r)
	go w.StartDispatch()
	go w.StartExecutor()

}
