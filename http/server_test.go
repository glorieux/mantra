package http_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"pkg.glorieux.io/mantra"
	mantra_http "pkg.glorieux.io/mantra/http"
	"pkg.glorieux.io/mantra/internal/log"
)

func TestHTTPServer(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	server := mantra_http.NewServer(routes)
	err := mantra.New(server)
	assert.NoError(t, err)

	server.Serve()

	res, err := http.Get(server.URL() + "/test")
	assert.NoError(t, err)
	respBody, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	assert.NoError(t, err)
	assert.Equal(t, "Response\n", string(respBody))
	assert.NoError(t, server.Stop())
}

func routes(router *mux.Router) {
	router.HandleFunc("/test", testHandler)
}

func testHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "Response")
}
