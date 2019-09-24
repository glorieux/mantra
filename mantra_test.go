package mantra_test

import (
	"context"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"pkg.glorieux.io/mantra"
)

var logger = logrus.New()

func TestMain(m *testing.M) {
	logger.Out = os.Stdout
	logger.Level = logrus.DebugLevel
	os.Exit(m.Run())
}

type testService struct {
	name     string
	Messages []string
}

func (ts *testService) Serve(ctx context.Context, app mantra.Application) error {
	_, mailbox := app.NewMailbox(ts.name)
	mailbox.Receive(func(message interface{}) {
		switch message.(type) {
		case testMessage:
			ts.Messages = append(ts.Messages, message.(testMessage).msg)
			message.(testMessage).ack.Send(true)
		default:
			ts.Messages = append(ts.Messages, message.(string))
		}
	})
	return nil
}

func (ts *testService) Stop() error {
	return nil
}

func (ts *testService) String() string {
	return ts.name
}

type testMessage struct {
	msg string
	ack *mantra.Address
}

func (testMessage) To() string { return "ts2" }

func TestServiceCommunication(t *testing.T) {
	ts1 := &testService{name: "ts1", Messages: []string{}}
	ts2 := &testService{name: "ts2", Messages: []string{}}
	app, err := mantra.New(logger, ts1, ts2)
	if err != nil {
		t.Error(err)
	}

	// ts1Address := app.Lookup("ts1")
	ts2Address := app.Lookup("ts2")

	t.Run("Sends message", func(t *testing.T) {
		ts2Address.Send(testMessage{"hello", app.Lookup("ts1")})

		if len(ts2.Messages) < 1 {
			t.Error("Should have received at least one message")
			return
		}

		if ts2.Messages[0] != "hello" {
			t.Error("Wrong message")
		}
	})

	// t.Run("Remove service", func(t *testing.T) {
	// 	toBeRemovedService := &testService{}
	// 	err := ts1Address.Send(mantra.AddServiceMessage(toBeRemovedService))
	// 	if err != nil {
	// 		t.Error(err)
	// 		return
	// 	}
	// 	err = ts1Address.send(mantra.RemoveServiceMessage(toBeRemovedService.String()))
	// 	if err != nil {
	// 		t.Error(err)
	// 		return
	// 	}
	// })
}
