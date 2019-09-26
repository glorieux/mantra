package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	"pkg.glorieux.io/mantra"
	"pkg.glorieux.io/mantra/example/greet"
	"pkg.glorieux.io/mantra/example/hello"
)

func main() {
	debug := flag.Bool("debug", false, "Print debug logs")
	flag.Parse()

	log := logrus.New()
	if *debug {
		log.SetLevel(logrus.DebugLevel)
	}

	hello := &hello.Service{}

	_, err := mantra.New(log, hello, &greet.Service{})
	if err != nil {
		log.Fatal(err)
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
