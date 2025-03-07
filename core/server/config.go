package server

import "net"

type Config struct {
	host string
	port int

	debug bool
}

func NewConfig() *Config {
	c := &Config{
		host:  "0.0.0.0",
		port:  9815,
		debug: false,
	}

	return c
}

func (c *Config) SetDebug(debug bool) *Config {
	c.debug = debug
	return c
}

func (c *Config) SetHost(host string) *Config {
	if net.ParseIP(host) == nil {
		c.host = host
	} else {
		c.host = "0.0.0.0"
	}
	return c
}

func (c *Config) SetPort(port int) *Config {
	if port < 0 || port > 65535 {
		c.port = 9815
	} else {
		c.port = port
	}
	return c
}
