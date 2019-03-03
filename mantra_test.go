package mantra_test

import (
	"context"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"glorieux.io/mantra"
)

var logger = logrus.New()

func TestMain(m *testing.M) {
	logger.Out = os.Stdout
	logger.Level = logrus.DebugLevel
	os.Exit(m.Run())
}

type testService struct {
	name     string
	send     mantra.SendFunc
	Messages []mantra.Message
}

func (ts *testService) Serve(ctx context.Context, msgChan <-chan mantra.Message, send mantra.SendFunc) error {
	ts.send = send
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

func (ts *testService) String() string {
	return ts.name
}

type unknownServiceMessage string

func (unknownServiceMessage) To() string { return "unknown_service" }

type testMessage struct {
	msg string
	ack chan bool
}

func (testMessage) To() string { return "ts2" }

func TestServiceCommunication(t *testing.T) {
	ts1 := &testService{name: "ts1", Messages: []mantra.Message{}}
	ts2 := &testService{name: "ts2", Messages: []mantra.Message{}}
	err := mantra.New(logger, ts1, ts2)
	if err != nil {
		t.Error(err)
	}
	t.Run("Unknow service", func(t *testing.T) {
		err := ts1.send(unknownServiceMessage("test"))
		if err == nil {
			t.Error("Should return Unknow service error")
		}
	})

	t.Run("Sends message", func(t *testing.T) {
		ack := make(chan bool)
		err := ts1.send(testMessage{"hello", ack})
		if err != nil {
			t.Error(err)
			return
		}
		<-ack

		if len(ts2.Messages) < 1 {
			t.Error("Should have received at least one message")
			return
		}

		if ts2.Messages[0].(testMessage).msg != "hello" {
			t.Error("Wrong message")
		}
	})

	t.Run("Remove service", func(t *testing.T) {
		toBeRemovedService := &testService{}
		err := ts1.send(mantra.AddServiceMessage{toBeRemovedService})
		if err != nil {
			t.Error(err)
			return
		}
		err = ts1.send(mantra.RemoveServiceMessage(toBeRemovedService.String()))
		if err != nil {
			t.Error(err)
			return
		}
	})
}
