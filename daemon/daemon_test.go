package daemon_test

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"pkg.glorieux.io/mantra/daemon"
)

var logger = logrus.New()

func TestMain(m *testing.M) {
	logger.Out = os.Stdout
	logger.Level = logrus.DebugLevel
	os.Exit(m.Run())
}

func TestStart(t *testing.T) {
	daemon.New()
}

// TODO
// There should be two sides of the daemon
// The first one (#1) is its definition and start. (The actual daemon)
// The second (#2) is the communication methods available from the outside.

// #1 daemon.New() <= Sets up a new Daemon
// #2 daemon.Start() <= Starts the daemon. It can be called multiple times.
