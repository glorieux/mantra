package mantra

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/thejerf/suture"
)

type MockedService struct {
	mock.Mock
}

func (m *MockedService) HandleMessage(message Message) error {
	return nil
}

func (m *MockedService) String() string {
	return "mockedService"
}

type message bool

func (message) To() string { return "test" }

func TestHandleMessage(t *testing.T) {
	serviceRegistry := newServiceRegistry(suture.NewSimple("test"), logrus.New())

	t.Run("AddServiceMessage", func(t *testing.T) {
		err := serviceRegistry.HandleMessage(AddServiceMessage{&MockedService{}})
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("RemoveServiceMessage", func(t *testing.T) {
		err := serviceRegistry.HandleMessage(RemoveServiceMessage("test"))
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("UnknownMessage", func(t *testing.T) {
		err := serviceRegistry.HandleMessage(message(true))
		if err == nil {
			t.Error("Should return an unknow message error")
		}
	})
}
