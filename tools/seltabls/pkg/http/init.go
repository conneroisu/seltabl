package http

import (
	"math/rand"
	"regexp"
	"time"
)

const (
	networkAdressPattern = `(?i)^((?:(?:tcp|udp|ip)[46]?)|(?:unix(?:gram|packet)?))://(.+)$`

	minTCPPort         = 0
	maxTCPPort         = 65535
	maxReservedTCPPort = 1024
	maxRandTCPPort     = maxTCPPort - (maxReservedTCPPort + 1)
)

var (
	tcpPortRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	homeDir     string
	homeDirSet  bool
	netAddrRx   *regexp.Regexp
)

// init initializes the ip
func init() {
	netAddrRx = regexp.MustCompile(networkAdressPattern)
}
