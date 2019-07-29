package dbops

import "time"

type MonitorRecord struct {
	Id         int64
	Type       string
	Value      string
	IP         string
	CreateTime time.Time
	Properties string
}
