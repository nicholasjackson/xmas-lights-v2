package main

import (
	"flag"
	"log"

	"github.com/matryer/vice/queues/sqs"
	"github.com/nicholasjackson/rcswitch"
	"github.com/nicholasjackson/xmas-lights-v2/lights"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/host"
)

var (
	pin      = flag.String("pin", "", "Name of the pin the 433MHz transmitter is attached to")
	sqsURI   = flag.String("sqs_uri", "", "URI for SQS queue to listen to")
	onCode   = flag.String("on_code", "", "Code to turn the switch on")
	offCode  = flag.String("off_code", "", "Code to turn the switch off")
	protocol = flag.Int("protocol", 0, "Switch protocol to use")
)

func main() {
	flag.Parse()
	log.Println("Starting lights: v0.1")

	_, err := host.Init()
	if err != nil {
		log.Fatal(err)
	}

	// Setup the GPIO settings for the 433MHz transmitter
	p := gpioreg.ByName(*pin)
	if p == nil {
		log.Fatal("Unable to find pin: ", *pin)
	}

	p.Out(gpio.High)
	log.Printf("%s: %s\n", p, p.Function())

	sw := rcswitch.New(p)

	c := lights.Config{
		OnCode:   *onCode,
		OffCode:  *offCode,
		Protocol: *protocol,
		SqsURI:   *sqsURI,
	}

	// Create a new SQS Vice transport
	t := sqs.New()

	// Create a new instances of the lights application and setup the transport
	l := lights.New(c, t, sw)
	err = l.Setup()
	if err != nil {
		log.Panicln(err)
	}

	log.Println("Listening for messages")
	l.Listen()
}
