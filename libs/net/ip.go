package net

import (
	"bytes"
	"errors"
	"fmt"
	"net"
)

// 通过提供子网段的方式获得ip
func LocalIpByIpNet(ipNet net.IPNet) (ip string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil && bytes.Equal(ipnet.Mask, ipNet.Mask) {
				ip = ipnet.IP.String()
			}
		}
	}
	if ip == "" {
		err = errors.New("can't find this subnet")
	}
	return
}
