package mantra_test

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"techmantra.io/mantra"
)

var logger = logrus.New()

func TestMain(m *testing.M) {
	logger.Out = os.Stdout
	logger.Level = logrus.DebugLevel
	os.Exit(m.Run())
}

type testApplication struct {
	send mantra.SendFunc
}

func (t *testApplication) Init(send mantra.SendFunc) error {
	t.send = send
	return nil
}

func (*testApplication) String() string {
	return "test"
}

type testService struct {
	messages []mantra.Message
}

func (ts *testService) HandleMessage(message mantra.Message) error {
	switch message.(type) {
	case testMessage:
		ts.messages = append(ts.messages, message)
		message.(testMessage).ack <- true
	default:
		ts.messages = append(ts.messages, message)
	}

	return nil
}

func (*testService) String() string {
	return "test"
}

func TestNewApplication(t *testing.T) {
	err := mantra.New(&testApplication{}, logger)
	if err != nil {
		t.Error(err)
	}
}

type unknownServiceMessage string

func (unknownServiceMessage) To() string { return "unknown_service" }

type testMessage struct {
	msg string
	ack chan bool
}

func (testMessage) To() string { return "test" }

func TestServiceCommunication(t *testing.T) {
	ta := &testApplication{}
	mantra.New(ta, logger)

	service := &testService{}
	err := ta.send(mantra.AddServiceMessage{service})
	if err != nil {
		t.Error(err)
		return
	}

	t.Run("Unknow service", func(t *testing.T) {
		err := ta.send(unknownServiceMessage("test"))
		if err == nil {
			t.Error("Should return Unknow service error")
		}
	})

	t.Run("Sends message", func(t *testing.T) {

		ack := make(chan bool)
		err = ta.send(testMessage{"hello", ack})
		if err != nil {
			t.Error(err)
			return
		}
		<-ack

		if len(service.messages) < 1 {
			t.Error("Should have received at least one message")
			return
		}

		if service.messages[0].(testMessage).msg != "hello" {
			t.Error("Wrong message")
		}
	})

	t.Run("Remove service", func(t *testing.T) {
		toBeRemovedService := &testService{}
		err := ta.send(mantra.AddServiceMessage{toBeRemovedService})
		if err != nil {
			t.Error(err)
			return
		}
		err = ta.send(mantra.RemoveServiceMessage(toBeRemovedService.String()))
		if err != nil {
			t.Error(err)
			return
		}
	})
}
