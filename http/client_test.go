package http_test

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"pkg.glorieux.io/mantra"
	mantraHTTP "pkg.glorieux.io/mantra/http"
	"pkg.glorieux.io/mantra/internal/log"
)

func TestHTTPClient(t *testing.T) {
	log.SetLevel(logrus.DebugLevel)
	client := mantraHTTP.NewClient()
	err := mantra.New(client)
	assert.NoError(t, err)
}
