package http

import (
	"fmt"
	"net"
)

// GetLocalIP returns the non loopback local IP of the host.
// It checks the address type and if it is not a loopback then
// it returns the IP.
//
// Example:
//
//	ip, err := GetLocalIP()
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(ip)
//	// Output:
//	// 192.168.1.100
func GetLocalIP() (ip string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", fmt.Errorf("could not get interface addresses: %w", err)
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
				break
			}
		}
	}
	return ip, nil
}
