package mantra

import (
	"os"
)

type Address struct {
	// 0 for local node or
	hostname int
	pid      int
	name     int
}

func newAddress(hostname, name string) *Address {
	return &Address{
		hostname: hostname,
		pid:      os.Getpid(),
		name:     name,
	}
}
