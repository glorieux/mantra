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

func (m *MockedService) Serve(ctx context.Context, msgChan <-chan Message, send SendFunc) error {
	return nil
}

func (*MockedService) Stop() error {
	return nil
}

func (m *MockedService) String() string {
	return "mockedService"
}

type message bool

func (message) To() string { return "test" }

func TestHandleMessage(t *testing.T) {
	supervisor := suture.NewSimple("test")
	supervisor.ServeBackground()
	serviceRegistry := newServiceRegistry(supervisor, logrus.New())

	t.Run("AddServiceMessage", func(t *testing.T) {
		err := serviceRegistry.send(AddServiceMessage{&MockedService{}})
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("RemoveServiceMessage", func(t *testing.T) {
		err := serviceRegistry.send(RemoveServiceMessage("test"))
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("UnknownMessage", func(t *testing.T) {
		err := serviceRegistry.send(message(true))
		if err == nil {
			t.Error("Should return an unknow message error")
		}
	})
}
