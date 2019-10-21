package daemon

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/spf13/afero"
	"pkg.glorieux.io/mantra"
	"pkg.glorieux.io/mantra/internal/log"
)

const (
	name = "mantra_daemon"
)

var (
	fs   afero.Fs
	once sync.Once

	config Config
)

// Daemon is a process running in the background
type Daemon interface {
	mantra.Service
}

type Config struct {
	Command string

	cachePath string
}

func binaryName(command string) string {
	return filepath.Base(command)
}

func (c Config) pidFileName() string {
	return fmt.Sprintf("%s.pid", binaryName(c.Command))
}

func (c Config) socketFileName() string {
	return fmt.Sprintf("%s.sock", binaryName(c.Command))
}

func init() {
	fs = afero.NewOsFs()
	cacheDirectory, err := os.UserCacheDir()
	if err != nil {
		log.Fatal(err)
	}
	config.cachePath = filepath.Join(cacheDirectory, name)
	err = os.MkdirAll(config.cachePath, 0700)
	if err != nil {
		log.Fatal(err)
	}
}

// Start starts a new daemon
func Start(c Config) {
	config.Command = c.Command
}

// New defines a new daemon
func New(debug bool, services ...mantra.Service) {
	if debug {
		log.SetLevel(log.DebugLevel)
	}
	file, err := os.OpenFile(
		filepath.Join(config.cachePath, fmt.Sprintf("%s.log", name)),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0666,
	)
	if err != nil {
		log.Warn("Failed to log to file, using default stdout", err)
	} else {
		log.Infof("Logging into: %s", config.cachePath)
		log.SetOutput(file)
	}

	err = mantra.New(services...)
	if err != nil {
		log.Fatal("Mantra error", err)
	}

	config.Command, err = os.Executable()
	if err != nil {
		log.Fatal("Mantra error", err)
	}
	log.Info("config", config)

	err = ioutil.WriteFile(
		filepath.Join(config.cachePath, config.pidFileName()),
		[]byte(strconv.Itoa(os.Getpid())),
		0666,
	)
	if err != nil {
		log.Error(err)
	}
}

func Stop() {
	d := getDaemon()
	// TODO: Send stop message
	d.conn.Close()
	d.process.Signal(os.Interrupt)
}
