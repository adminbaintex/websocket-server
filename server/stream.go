package server

import (
	"errors"
	"github.com/gorilla/websocket"
)

// Stream interface simplify the access to a connection receiving messages
// from channel and encapsulation messages in a convenient struct
type Stream interface {
	Incoming() chan *WSMessage
	Send(*WSMessage) error
	Close()
}

// AbstractStream implements the communication with a websocket connection
type AbstractStream struct {
	incomingChan chan *WSMessage

	conn *websocket.Conn

	done chan bool
	quit chan bool
}

func newAbstractStream(conn *websocket.Conn) Stream {
	as := &AbstractStream{
		incomingChan: make(chan *WSMessage),
		done:         make(chan bool),
		quit:         make(chan bool),
		conn:         conn,
	}
	go as.read()

	return as
}

func (as *AbstractStream) read() {

	defer func() {
		close(as.incomingChan)
		as.done <- true
	}()

	for {

		select {
		case <-as.quit:
			return
		default:
		}

		mt, data, err := as.conn.ReadMessage()
		if err != nil {
			// Not logged: this fires on every client disconnect (read timeout,
			// connection closed, ...) and the handler already reports it.
			return
		}
		message := new(WSMessage)
		message.Data = data
		message.MessageType = mt

		// Deliver without wedging: incomingChan is unbuffered, so if the
		// consumer stopped reading (it decided to disconnect while this message
		// was in flight) a bare send would park this goroutine forever and
		// Close() — which waits for it via <-as.done — would deadlock.
		select {
		case as.incomingChan <- message:
		case <-as.quit:
			return
		}
	}
}

// Incoming returns the channel used for reading incoming packets.
// The channel gets automatically closed when the stream gets closed.
func (as *AbstractStream) Incoming() chan *WSMessage {
	return as.incomingChan
}

// Send will safely write the message
func (as *AbstractStream) Send(m *WSMessage) error {
	if m == nil {
		return errors.New("WSMessage parameter is nil")
	}
	return as.conn.WriteMessage(m.MessageType, m.Data)
}

// Close will close the stream and cleanup open channels and running
// go routines.
func (as *AbstractStream) Close() {
	as.conn.Close()
	close(as.quit)
	<-as.done
}
