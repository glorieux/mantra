package mantra_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"pkg.glorieux.io/mantra"
	"pkg.glorieux.io/mantra/internal/log"
)

func TestMain(m *testing.M) {
	log.SetLevel(log.DebugLevel)
	os.Exit(m.Run())
}

type testService struct{}

func (*testService) TestWithoutParams() {}

func (*testService) TestMessage(m chan string) {
	m <- "test"
}

func (ts *testService) Stop() error {
	return nil
}

func TestServiceCommunication(t *testing.T) {
	ts1 := &testService{}

	err := mantra.New(ts1)
	if err != nil {
		t.Error(err)
	}

	ack := make(chan string)
	mantra.Send(mantra.Lookup("testService"), ts1.TestWithoutParams)
	mantra.Send(mantra.Lookup("testService"), ts1.TestMessage, ack)
	assert.NotEmpty(t, <-ack)
	mantra.Stop()
}
