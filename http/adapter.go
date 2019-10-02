package http

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/pprof"
	"strings"
	"time"

	oidc "github.com/coreos/go-oidc"
	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
)

// TODO: Rename to Middleware and integrate to HTTP Server

// Adapter enables http middlewares
type Adapter func(http.Handler) http.Handler

// Adapt applies the given adapters
func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

// BasicAuthentication adds basic HTTP authentication
func BasicAuthentication(username, password, realm string) Adapter {
	checkAuth := func(r *http.Request) string {
		s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(s) != 2 || s[0] != "Basic" {
			return ""
		}
		b, err := base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			return ""
		}
		pair := strings.SplitN(string(b), ":", 2)
		if len(pair) != 2 {
			return ""
		}
		if pair[0] == username && pair[1] == password {
			return pair[0]
		}
		return ""
	}

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u := checkAuth(r)
			if u == "" {
				w.Header().Set("Content-Type", "text/plain")
				w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, realm))
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprint(w, http.StatusText(http.StatusUnauthorized))
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}

// CloudflareAuthentication verifies cloudflare authentication tokens
func CloudflareAuthentication(authDomain, policyAUD string) Adapter {
	const jwtHeaderKey = "Cf-Access-Jwt-Assertion"

	var (
		certsURL = fmt.Sprintf("%s/cdn-cgi/access/certs", authDomain)
		keySet   = oidc.NewRemoteKeySet(context.Background(), certsURL)
		verifier = oidc.NewVerifier(authDomain, keySet, &oidc.Config{ClientID: policyAUD})
	)

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accessJWT := r.Header.Get(jwtHeaderKey)
			if accessJWT == "" {
				http.Error(w, "No token on the request", http.StatusUnauthorized)
				return
			}

			_, err := verifier.Verify(r.Context(), accessJWT)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}

// Logging logs the request basic informations
func Logging(l *log.Logger) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l.Infof("%s %s", r.Method, r.URL.Path)
			h.ServeHTTP(w, r)
		})
	}
}

// Header adds a given header to the response
func Header(key, value string) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add(key, value)
			h.ServeHTTP(w, r)
		})
	}
}

// Cache adds cache-control headers to the response
func Cache() Adapter {
	return Header("Cache-Control", "public, max-age=2592000")
}

// Timing adds request timing information
func Timing(l *log.Logger) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			beginning := time.Now()
			h.ServeHTTP(w, r)
			end := time.Now()
			l.Infof("%s %s took %v to run", r.Method, r.URL.Path, end.Sub(beginning))
		})
	}
}

// TurbolinksRedirect handles turbolinks redirections
func TurbolinksRedirect(store sessions.Store, sessionName string, location string) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := store.Get(r, sessionName)
			if err == nil {
				turbolinksLocation := session.Values[location]
				if turbolinksLocation != nil {
					w.Header().Add("Turbolinks-Location", turbolinksLocation.(string))
					session.Options.MaxAge = -1
					session.Save(r, w)
				}
			}
			h.ServeHTTP(w, r)
		})
	}
}

// PProf adds pprof handlers
func PProf() Adapter {
	return func(h http.Handler) http.Handler {
		log.Warn("pprof is enabled. It might create a bit of overhead.")
		mux := http.NewServeMux()
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
		mux.Handle("/", h)
		return mux
	}
}
