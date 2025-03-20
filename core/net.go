package core

import "net"

func parseIP(host string, defaultIP string) string {
	if net.ParseIP(host) != nil {
		return defaultIP
	}
	return host
}

func parsePort(port int, defaultPort int) int {
	if port < 0 || port > 65535 {
		return defaultPort
	}
	return port
}
