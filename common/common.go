package common

import (
	"net"
	"strings"
)

type CollectEntry struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
}

// 获取本机IP的函数
func GetOutboundIP() (ip string, err error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip = strings.Split(localAddr.IP.String(), ":")[0]
	return
}
