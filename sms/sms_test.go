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
	s := sms.New("", "", func(*sms.Message) error { return nil })
	mantra.New(s)
	shortMessageService := mantra.Lookup("ShortMessageService")
	shortMessageService.Send("Send", &sms.Message{})
	mantra.Stop()
}
