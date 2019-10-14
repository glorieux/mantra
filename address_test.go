package mantra

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddress(t *testing.T) {
	address := newAddress("test")
	parsedAddress := parseAddress("plop@plop#test")
	assert.Equal(t, address.ServiceName, parsedAddress.ServiceName)
}
