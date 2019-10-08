package update

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"pkg.glorieux.io/mantra"
	"pkg.glorieux.io/mantra/internal/log"
	"pkg.glorieux.io/version"
)

// Service is a service handling self-updates
type Service struct {
	ticker   *time.Ticker
	provider Provider

	sync.Mutex
}

// Provider provides updates
type Provider interface {
	Interval() time.Duration
	Versions() ([]*version.Version, error)
	LatestVersion() (*version.Version, error)
	Download(version *version.Version) error
}

// New returns a new update service
func New(provider Provider) mantra.Service {
	return &Service{
		ticker:   time.NewTicker(provider.Interval()),
		provider: provider,
	}
}

// Serve runs the service
// TODO also handle manual version check
func (s *Service) Serve(mux mantra.ServeMux) {
	s.check()
	go func() {
		for {
			<-s.ticker.C
			s.check()
		}
	}()
}

// Stop stops the service
func (s *Service) Stop() error {
	s.ticker.Stop()
	return nil
}

func (*Service) String() string {
	return "update"
}

func (s *Service) check() {
	s.Lock()
	defer s.Unlock()
	log.Info("Checking for updates")
	versions, err := s.provider.Versions()
	if err != nil {
		logrus.Error("Could not check available versions", err)
		return
	}

	// TODO if current remote version > current version
	log.Info("Versions: ", versions)
}
