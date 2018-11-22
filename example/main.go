package main

import (
	"flag"
	"os"

	"github.com/sirupsen/logrus"
	"techmantra.io/mantra"
	"techmantra.io/mantra/example/greet"
	"techmantra.io/mantra/example/hello"
)

func main() {
	debug := flag.Bool("debug", false, "Print debug logs")
	flag.Parse()

	log := logrus.New()
	if *debug {
		log.SetLevel(logrus.DebugLevel)
	}
	r := &greeter{}
	err := mantra.New(r, log)
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}
	greeted := make(chan bool)
	r.send(greet.Message{"world", greeted})
	<-greeted
}

type greeter struct {
	send mantra.SendFunc
}

func (g *greeter) Init(send mantra.SendFunc) error {
	g.send = send
	send(mantra.AddServiceMessage{&greet.Service{}})
	send(mantra.AddServiceMessage{&hello.Service{}})
	return nil
}
func (*greeter) String() string {
	return "greeter"
}
