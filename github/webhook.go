package github

import (
	"io/ioutil"
	"log"
	stdHttp "net/http"

	"github.com/gorilla/mux"
	"pkg.glorieux.io/mantra"
	"pkg.glorieux.io/mantra/http"
)

func NewWebhookServer(msgChan chan string) mantra.Service {
	return http.NewServer(func(router *mux.Router) {
		router.HandleFunc("/hook", func(w stdHttp.ResponseWriter, r *stdHttp.Request) {
			respBody, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Fatal(err)
			}
			err = r.Body.Close()
			if err != nil {
				log.Fatal(err)
			}
			msgChan <- string(respBody)
		})
	})
}
