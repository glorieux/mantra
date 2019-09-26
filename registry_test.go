package mantra

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/thejerf/suture"
)

type MockedService struct {
	mock.Mock
}

func (m *MockedService) Serve(ctx context.Context, app Application) error {
	return nil
}

func (*MockedService) Stop() error {
	return nil
}

func (m *MockedService) String() string {
	return "mockedService"
}

type testMessage bool

func (message) To() string { return "test" }

func TestHandleMessage(t *testing.T) {
	supervisor := suture.NewSimple("test")
	supervisor.ServeBackground()
	serviceRegistry := newServiceRegistry(supervisor, logrus.New())
	serviceRegistryAddress := serviceRegistry.Lookup(registryServiceName)

	t.Run("AddServiceMessage", func(t *testing.T) {
		serviceRegistryAddress.Send(AddServiceMessage(&MockedService{}))
	})

	t.Run("RemoveServiceMessage", func(t *testing.T) {
		serviceRegistryAddress.Send(RemoveServiceMessage("test"))
	})

	t.Run("UnknownMessage", func(t *testing.T) {
		serviceRegistryAddress.Send(true)
	})
}
