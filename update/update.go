package update

import (
	"sort"
	"time"

	"github.com/sirupsen/logrus"
	"pkg.glorieux.io/mantra"
	"pkg.glorieux.io/mantra/internal/log"
	"pkg.glorieux.io/version"
)

// UpdateService is a service handling self-updates
type UpdateService struct {
	current  *version.Version
	provider Provider
}

// Provider provides updates
type Provider interface {
	Interval() time.Duration
	Versions() ([]*version.Version, error)
	Download(version *version.Version) error
}

// New returns a new update service
func New(current *version.Version, provider Provider) mantra.Service {
	s := &UpdateService{
		current:  current,
		provider: provider,
	}
	go s.serve()
	return s
}

// Serve runs the service
// TODO also handle manual version check
func (s *UpdateService) serve() {
	for {
		time.Sleep(s.provider.Interval())
		s.check()
	}
}

// Stop stops the service
func (s *UpdateService) Stop() error {
	return nil
}

func (s *UpdateService) check() {
	log.Info("Checking for updates")
	versions, err := s.provider.Versions()
	if err != nil {
		logrus.Error("Could not check available versions", err)
		return
	}
	sort.Sort(version.Ascending(versions))
	lastVersion := versions[len(versions)-1]

	if s.current.Equal(lastVersion) {
		log.Debug("Equal")
		return
	}
	if s.current.After(lastVersion) {
		log.Debug("After")
		return
	}
	err = s.provider.Download(lastVersion)
	if err != nil {
		log.Error(err)
	}

	// TODO
	// * Replace binary
	// * Run healcheck
	// * Send stop message
	mantra.Stop()
	return
}
