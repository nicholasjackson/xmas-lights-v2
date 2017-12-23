package lights

import (
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/nicholasjackson/rcswitch"
)

var config = Config{
	SqsURI: "https://something.com/something",
}

func setup(t *testing.T, rc <-chan []byte) (*is.I, *Lights, *ViceTransportMock, *rcswitch.SwitchMock) {
	mockedViceTransport := &ViceTransportMock{
		DoneFunc: func() chan struct{} {
			panic("TODO: mock out the Done method")
		},
		ErrChanFunc: func() <-chan error {
			panic("TODO: mock out the ErrChan method")
		},
		ReceiveFunc: func(name string) <-chan []byte {
			//			panic("TODO: mock out the Receive method")
			return rc
		},
		SendFunc: func(name string) chan<- []byte {
			panic("TODO: mock out the Send method")
		},
		StopFunc: func() {
			panic("TODO: mock out the Stop method")
		},
	}

	mockedSwitch := &rcswitch.SwitchMock{
		ScanFunc: func() {
			panic("TODO: mock out the Scan method")
		},
		SendFunc: func(data string, protocolID int) {
		},
	}

	is := is.New(t)

	l := &Lights{
		config:    config,
		transport: mockedViceTransport,
		swtch:     mockedSwitch,
	}

	return is, l, mockedViceTransport, mockedSwitch
}

func TestConfiguresTransportOnSetup(t *testing.T) {
	is, l, mt, _ := setup(t, nil)
	l.Setup()

	is.Equal(1, len(mt.ReceiveCalls()))                // should have called receive once
	is.Equal(config.SqsURI, mt.ReceiveCalls()[0].Name) // should have called receive with correct sqs uri
}

func TestOnMessageSwitchesOn(t *testing.T) {
	mc := make(chan []byte)
	is, l, _, ms := setup(t, mc)
	l.Setup()

	sendAndWait(l, mc, `{"command":"on"}`)

	is.Equal(1, len(ms.SendCalls()))                        // Should have called send once
	is.Equal(config.OnCode, ms.SendCalls()[0].Data)         // Should have called with the correct code
	is.Equal(config.Protocol, ms.SendCalls()[0].ProtocolID) // Should have called with the correct protocolID
}

func TestOffMessageSwitchesOn(t *testing.T) {
	mc := make(chan []byte)
	is, l, _, ms := setup(t, mc)
	l.Setup()

	sendAndWait(l, mc, `{"command":"off"}`)

	is.Equal(1, len(ms.SendCalls()))                        // Should have called send once
	is.Equal(config.OffCode, ms.SendCalls()[0].Data)        // Should have called with the correct code
	is.Equal(config.Protocol, ms.SendCalls()[0].ProtocolID) // Should have called with the correct protocolID
}

func TestInvalidMessageDoesNothing(t *testing.T) {
	mc := make(chan []byte)
	is, l, _, ms := setup(t, mc)
	l.Setup()

	sendAndWait(l, mc, `{"command":"wibble"}`)

	is.Equal(0, len(ms.SendCalls())) // Should have called send 0 times
}

func sendAndWait(l *Lights, rc chan []byte, message string) {
	go func() {
		l.Listen()
	}()

	time.Sleep(10 * time.Millisecond) // I know this is super hacky

	go func() {
		rc <- []byte(message)
	}()

	time.Sleep(10 * time.Millisecond) // I know this is super hacky
}
