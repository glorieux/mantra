package mantra

import (
	"fmt"
	"strings"
)

type topic struct {
	Address *Address
	Method  string
}

func newTopic(address *Address, method string) *topic {
	return &topic{
		Address: address,
		Method:  method,
	}
}

func parseTopic(t string) *topic {
	return &topic{
		Address: parseAddress(t),
		Method:  t[strings.LastIndex(t, ".")+1:],
	}
}

func (t *topic) String() string {
	return fmt.Sprintf("%s.%s", t.Address.String(), t.Method)
}
