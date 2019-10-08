package mantra

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thejerf/suture"
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
	registry := newServiceRegistry(suture.NewSimple("test"))

	err := registry.addService(&MockedService{name: ""})
	if assert.Error(t, err) {
		assert.Equal(t, ErrEmptyName, err)
	}
	err = registry.addService(&MockedService{name: registryServiceName})
	if assert.Error(t, err) {
		assert.Equal(t, fmt.Errorf("Do not use %s as service name", registryServiceName), err)
	}

	err = registry.addService(&MockedService{name: "@#$%^&*("})
	if assert.Error(t, err) {
		assert.Equal(t, ErrNotAlphaNum, err)
	}
}
