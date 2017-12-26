package lights

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/matryer/vice"
	"github.com/nicholasjackson/rcswitch"
)

const (
	// CommandTurnOn turns the lights on
	CommandTurnOn = "on"
	// CommandTurnOff turns the lights off
	CommandTurnOff = "off" // Turn the ligths off
)

// Message defines the format of the message
type Message struct {
	Command string `json:"command"`
}

//go:generate moq -out transport_moq_test.go . ViceTransport

// ViceTransport wraps the vice.Transport interface to allow mock creation
type ViceTransport interface {
	vice.Transport
}

// Lights represents an instance of the lights swtiching app
type Lights struct {
	transport   vice.Transport
	swtch       rcswitch.Switch
	config      Config
	messageChan <-chan []byte
}

// Config represents the configuration options for the application
type Config struct {
	OnCode   string
	OffCode  string
	Protocol int
	SqsURI   string
}

// New creates a new instance of Lights
func New(c Config, t vice.Transport, swtch rcswitch.Switch) *Lights {
	return &Lights{transport: t, swtch: swtch, config: c}
}

// Setup sets up the transport
func (l *Lights) Setup() error {
	l.messageChan = l.transport.Receive(l.config.SqsURI)
	if l.messageChan == nil {
		return fmt.Errorf("Unable to create connection to SQS, have you set the AWS credentials")
	}

	return nil
}

// Listen for message and turn on and off the switch
func (l *Lights) Listen() {
	for m := range l.messageChan {
		log.Println("Received Message:", m)

		message := decodeMessage(m)
		switch message.Command {
		case CommandTurnOn:
			log.Println("Sending turn on with code:", l.config.OnCode)
			l.swtch.Send(l.config.OnCode, l.config.Protocol)
		case CommandTurnOff:
			log.Println("Sending turn off with code:", l.config.OffCode)
			l.swtch.Send(l.config.OffCode, l.config.Protocol)
		}
	}
}

func decodeMessage(data []byte) Message {
	var m Message
	json.Unmarshal(data, &m)

	return m
}
