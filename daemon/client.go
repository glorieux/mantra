package daemon

import (
	"fmt"
	"io/ioutil"
	"net/rpc"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/sirupsen/logrus"
)

var (
	runningDaemon *remoteDaemon
)

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
