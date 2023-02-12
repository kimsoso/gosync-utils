package snet

import (
	"net"
	"time"
)

const (
	Timeout   = time.Millisecond * 100
	LoopTimes = 3
)

func LoopConn(host string) (conn net.Conn, err error) {
	dialer := net.Dialer{Timeout: Timeout}
	for i := 0; i < LoopTimes; i++ {
		if conn, err = dialer.Dial("tcp", host); err == nil {
			return
		}
	}
	return
}

func TimeoutConn(host string) (conn net.Conn, err error) {
	dialer := net.Dialer{Timeout: Timeout}
	return dialer.Dial("tcp", host)
}
