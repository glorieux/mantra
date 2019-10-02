package mantra

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockedService struct {
	mock.Mock
	name string
}

func (m *MockedService) Serve(mux ServeMux) {}

func (*MockedService) Stop() error {
	return nil
}

func (m *MockedService) String() string {
	return m.name
}

func TestAddService(t *testing.T) {
	registry := &registry{}
	assert.Panics(t, func() {
		registry.addService(&MockedService{name: ""})
	}, "Empty name")
	assert.Panics(t, func() {
		registry.addService(&MockedService{name: "mantra"})
	}, "Registry name")
	assert.Panics(t, func() {
		registry.addService(&MockedService{name: "*&(#&@^"})
	}, "special characters")

}
