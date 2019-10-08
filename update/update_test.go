package update

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"pkg.glorieux.io/mantra"
	"pkg.glorieux.io/mantra/internal/log"
	"pkg.glorieux.io/version"
)

type mockProvider struct {
	mock.Mock
}

func (m *mockProvider) Interval() time.Duration {
	args := m.Called()
	return args.Get(0).(time.Duration)
}

func (m *mockProvider) Versions() ([]*version.Version, error) {
	args := m.Called()
	return args.Get(0).([]*version.Version), args.Error(1)
}

func (m *mockProvider) LatestVersion() (*version.Version, error) {
	args := m.Called()
	return args.Get(0).(*version.Version), args.Error(1)
}

func (m *mockProvider) Download(version *version.Version) error {
	args := m.Called()
	return args.Error(0)
}

func TestMain(m *testing.M) {
	log.SetLevel(log.DebugLevel)
	os.Exit(m.Run())
}

func TestService(t *testing.T) {
	provider := &mockProvider{}
	provider.On("Interval").Return(5 * time.Second)

	v, _ := version.New("0.1.0")
	provider.On("Versions").Return([]*version.Version{v}, nil)

	mantra.New(New(provider))
	provider.AssertCalled(t, "Interval")
	provider.AssertCalled(t, "Versions")
}
