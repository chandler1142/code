package net_monitor

import (
	"sync"
)

type Iface struct {
	name     string
	ip       string
	lastSend uint64
	lastRecv uint64
	mtx      *sync.Mutex
}

func NewIfmon(name string, ip string) *Iface {
	return &Iface{
		name: name,
		ip:   ip,
		mtx:  &sync.Mutex{},
	}
}
