package github

import (
	"pkg.glorieux.io/mantra"
	"pkg.glorieux.io/mantra/http"
)

func NewGithubService() mantra.Service {
	client := http.NewClient()
	return client
}
