package utils

import (
	"net"
)

func GetIP() string {

	interfaceAddr, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}

	ips := make([]string, 0)
	for _, address := range interfaceAddr {
		ipNet, isValidIpNet := address.(*net.IPNet)
		if isValidIpNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ips = append(ips, ipNet.IP.String())
			}
		}
	}
	if len(ips) <= 0 {
		return "127.0.0.1"
	}
	return ips[0]
}
