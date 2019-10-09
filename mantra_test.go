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

type testService struct {
	name string
}

func (ts *testService) Receive(mux mantra.ServeMux) {
	mux.Handle("test", func(e mantra.Event) {
		e.Data.(chan string) <- "test"
	})
}

func (ts *testService) Serve() {}

func (ts *testService) Stop() error {
	return nil
}

func (ts *testService) String() string {
	return ts.name
}

func TestServiceCommunication(t *testing.T) {
	ts1 := &testService{name: "ts1"}
	err := mantra.New(ts1)
	if err != nil {
		t.Error(err)
	}

	ack := make(chan string)
	err = mantra.Send("ts1.test", ack)
	assert.NoError(t, err)
	assert.NotEmpty(t, <-ack)
	mantra.Stop()
}
