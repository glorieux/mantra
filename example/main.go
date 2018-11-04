package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"techmantra.io/mantra"
	"techmantra.io/mantra/example/greet"
	"techmantra.io/mantra/example/hello"
)

func main() {
	r := &greeter{}
	log := logrus.New()
	err := mantra.New(r, log)
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}
	greeted := make(chan bool)
	err = r.send(greet.Message{"world", greeted})
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}
	<-greeted
}

type greeter struct {
	send mantra.SendFunc
}

func (*greeter) String() string {
	return "greeter"
}

func (g *greeter) Init(send mantra.SendFunc) error {
	g.send = send
	send(mantra.AddServiceMessage{greet.New(send)})
	send(mantra.AddServiceMessage{&hello.Service{}})
	return nil
}
