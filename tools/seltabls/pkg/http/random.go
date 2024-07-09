package http

import (
	"fmt"
	"net"
)

// IsTCPPortAvailable returns a flag indicating whether or not a TCP port is
// available.
//
// It takes a port as an integer as an argument and returns a boolean indicating
// whether or not the port is available.
//
// Example:
//
//	package main
//
//	import "fmt"
//	import "github.com/cseval/cse-ncaa/pkg/ip"
//
//	func main() {
//		mux := http.NewServeMux()
//		mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
//			w.Write([]byte("Hello World"))
//		})
//		srv := &http.Server{
//			Addr: ":8080",
//			Handler: mux,
//		}
//		if err := srv.ListenAndServe(); err != nil {
//			log.Fatal(err)
//		}
//		if IsTCPPortAvailable(8080) {
//			fmt.Println("Port is available")
//		} else {
//			fmt.Println("Port is not available")
//		}
//	}
func IsTCPPortAvailable(port int) bool {
	if port < minTCPPort || port > maxTCPPort {
		return false
	}
	conn, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// RandomTCPPort gets a free, random TCP port between 1025-65535. If no free
// ports are available -1 is returned.
//
// It returns an integer representing a random TCP port between 1025-65535.
// If no free ports are available, it returns -1.
//
// Example:
//
//	package main
//
//	import "fmt"
//	import "github.com/cseval/cse-ncaa/pkg/ip"
//
//	func main() {
//		port := ip.RandomTCPPort()
//		fmt.Println(port)
//	}
func RandomTCPPort() int {
	for i := maxReservedTCPPort; i < maxTCPPort; i++ {
		p := tcpPortRand.Intn(maxRandTCPPort) + maxReservedTCPPort + 1
		if IsTCPPortAvailable(p) {
			return p
		}
	}
	return -1
}
