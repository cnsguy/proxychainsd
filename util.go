package main

import (
	"net"
	"time"
)

func RefreshTimeout(c net.Conn, seconds time.Duration) {
	c.SetDeadline(time.Now().Add(seconds * time.Second))
}
