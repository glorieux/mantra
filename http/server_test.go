package http_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"pkg.glorieux.io/mantra"
	mantra_http "pkg.glorieux.io/mantra/http"
)

func TestHTTPServer(t *testing.T) {
	var logger = logrus.New()
	logger.Out = os.Stdout
	logger.Level = logrus.DebugLevel

	server := mantra_http.NewServer(routes)
	err := mantra.New(logger, server)
	assert.NoError(t, err)

	res, err := http.Get(server.URL() + "/test")
	assert.NoError(t, err)
	respBody, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.NoError(t, err)
	assert.Equal(t, "Response\n", string(respBody))
}

func routes(router *mux.Router) {
	router.HandleFunc("/test", testHandler)
}

func testHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "Response")
}
