package net

import (
	"fmt"
	"net"
)

// ListenerFromPort returns a listener on the next available port starting from
// the given base port
func ListenerFromPort(base int) net.Listener {
	port := base

	for {
		ln, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
		if err == nil {
			return ln
		}
		port++
	}
}
