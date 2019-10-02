package http_test

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"pkg.glorieux.io/mantra"
	mantraHTTP "pkg.glorieux.io/mantra/http"
)

func TestHTTPClient(t *testing.T) {
	var logger = logrus.New()
	logger.Out = os.Stdout
	logger.Level = logrus.DebugLevel
	client := mantraHTTP.NewClient()
	err := mantra.New(logger, client)
	assert.NoError(t, err)
}
