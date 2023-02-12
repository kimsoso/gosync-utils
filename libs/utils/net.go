package utils

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

const (
	ALIDNS = "223.5.5.5:53"
)

func GetOutBoundIP() (ip string, err error) {
	conn, err := net.Dial("udp", ALIDNS)
	if err != nil {
		fmt.Println(err)
		return
	}

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	fmt.Println(localAddr.String())

	ip = strings.Split(localAddr.String(), ":")[0]

	return
}

func GetLocalIp() (ip string, err error) {
	ip, err = GetOutBoundIP()
	if err == nil {
		return
	}

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
				return
			}
		}
	}
	return "", errors.New("can't get local ip")
}
