package common

import (
	"net"
	"strings"
)

//要收集的日志的配置结构体
type CollectEntry struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
}

//获取ip地址
func GetOutIp() (ip string, err error) {
	con, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return
	}
	defer con.Close()
	localAddr := con.LocalAddr().(*net.UDPAddr)
	ip = strings.Split(localAddr.IP.String(), ":")[0]
	return
}
