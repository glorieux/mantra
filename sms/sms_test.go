package sms_test

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"pkg.glorieux.io/mantra"
	"pkg.glorieux.io/mantra/internal/log"
	"pkg.glorieux.io/mantra/sms"
)

func TestMain(m *testing.M) {
	log.SetLevel(logrus.DebugLevel)
	os.Exit(m.Run())
}

func TestSend(t *testing.T) {
	s := sms.New()
	mantra.New(s)
	mantra.SendMessage(mantra.Lookup("ShortMessageService"), s.Send, &sms.Message{})
	mantra.Stop()
}
