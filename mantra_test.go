package mantra_test

import (
	"context"
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
	Messages []mantra.Message
}

func (ts *testService) Serve(ctx context.Context, msgChan <-chan mantra.Message, send mantra.SendFunc) error {
	for message := range msgChan {
		switch message.(type) {
		case testMessage:
			ts.Messages = append(ts.Messages, message)
			close(message.(testMessage).ack)
		default:
			ts.Messages = append(ts.Messages, message)
		}
	}
	return nil
}

func (ts *testService) Stop() error {
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
	err := mantra.New(ta, logger)
	if err != nil {
		t.Error(err)
	}
	ts := &testService{Messages: []mantra.Message{}}
	err = ta.send(mantra.AddServiceMessage{ts})
	if err != nil {
		t.Error(err)
	}

	t.Run("Unknow service", func(t *testing.T) {
		err := ta.send(unknownServiceMessage("test"))
		if err == nil {
			t.Error("Should return Unknow service error")
		}
	})

	t.Run("Sends message", func(t *testing.T) {
		ack := make(chan bool)
		err := ta.send(testMessage{"hello", ack})
		if err != nil {
			t.Error(err)
			return
		}
		<-ack

		if len(ts.Messages) < 1 {
			t.Error("Should have received at least one message")
			return
		}

		if ts.Messages[0].(testMessage).msg != "hello" {
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
