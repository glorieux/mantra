package daemon

import (
	"fmt"
	"io/ioutil"
	"net/rpc"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"pkg.glorieux.io/mantra"
	"pkg.glorieux.io/mantra/internal/log"
)

var (
	fs   afero.Fs
	once sync.Once

	config        Config
	runningDaemon *remoteDaemon
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
	config.cachePath = filepath.Join(cacheDirectory, "mantra_daemon")
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
		filepath.Join(config.cachePath, "mantra_daemon.log"),
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

func getDaemon() *remoteDaemon {
	// TODO: Check for PID file and estabilshes connection through socket
	// TODO: Clear the PID on exit
	once.Do(func() {
		content, err := ioutil.ReadFile(filepath.Join(config.cachePath, config.pidFileName()))
		if err != nil || len(content) == 0 {
			fmt.Println("Command", config.Command, config)
			command := exec.Command(config.Command)
			command.Start()
		}
		content, err = ioutil.ReadFile(filepath.Join(config.cachePath, config.pidFileName()))
		if err != nil {
			logrus.Fatal(err)
		}

		pid, err := strconv.Atoi(string(content))
		if err != nil {
			logrus.Fatal(err)
		}

		process, err := os.FindProcess(pid)
		if err != nil {
			logrus.Fatal(err)
		}

		conn, err := rpc.Dial("unix", filepath.Join(config.cachePath, config.socketFileName()))
		if err != nil {
			logrus.Fatal(err)
		}
		runningDaemon = &remoteDaemon{
			process: process,
			conn:    conn,
		}
	})
	return runningDaemon
}

type remoteDaemon struct {
	process *os.Process
	conn    *rpc.Client
}

func TestConn() error {
	d := getDaemon()
	req := 42
	var reply int
	return d.conn.Call("daemon.ConnTest", &req, &reply)
}

func Stop() {
	d := getDaemon()
	// TODO: Send stop message
	d.conn.Close()
	d.process.Signal(os.Interrupt)
}
