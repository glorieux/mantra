package mantra_test

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"pkg.glorieux.io/mantra"
)

var logger = logrus.New()

func TestMain(m *testing.M) {
	logger.Out = os.Stdout
	logger.Level = logrus.DebugLevel
	os.Exit(m.Run())
}

type testService struct {
	name string
}

func (ts *testService) Serve(mux mantra.ServeMux) {
	mux.Handle("test", func(e mantra.Event) {
		e.Data.(chan string) <- "test"
	})
}

func (ts *testService) Stop() error {
	return nil
}

func (ts *testService) String() string {
	return ts.name
}

func TestServiceCommunication(t *testing.T) {
	ts1 := &testService{name: "ts1"}
	err := mantra.New(logger, ts1)
	if err != nil {
		t.Error(err)
	}

	ack := make(chan string)
	err = mantra.Send("ts1.test", ack)
	assert.NoError(t, err)
	assert.NotEmpty(t, <-ack)
}
