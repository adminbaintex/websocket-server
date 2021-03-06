package server

import (
	"github.com/gorilla/websocket"
	"log"
	"net"
	"net/http"
	"net/url"
	"github.com/armon/go-proxyproto"
)

// WSHandler will receive new connections as streams.
type WSHandler interface {
	ServeWS(net.Conn, Stream, url.Values)
}

// Server manages multiple Configurations and yields new connection as
// streams to the Handler.
type Server struct {
	// The Handler that receives new Streams.
	wsHandler WSHandler

	// Server buffers
	readBufferSize int
	writeBufferSize int

	// Server Listener and Handler
	listener    net.Listener
	httpHandler http.Handler

	// Use Proxy protocol
	ProxyProtocol bool
}

// NewServer returns a new Server.
func NewServer(handler WSHandler, readBufferSize int, writeBufferSize int, proxyProtocol bool) *Server {
	return &Server{
		wsHandler: handler,
		readBufferSize: readBufferSize,
		writeBufferSize: writeBufferSize,
		ProxyProtocol: proxyProtocol,
	}
}

func (s *Server) serve() {
	httpServer := &http.Server{Handler: s.httpHandler}
	err := httpServer.Serve(s.listener)
	if err != nil {
		log.Println(err)
	}
}

// ListenAndServe will run a simple WS (HTTP) server.
// return listening error
func (s *Server) ListenAndServe(address string, path string) error {
	l, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	if s.ProxyProtocol {
		// Wrap listener in a proxyproto listener
		s.listener = &proxyproto.Listener{Listener: l}
	} else {
		s.listener = l
	}

	mux := http.NewServeMux()
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  s.readBufferSize,
			WriteBufferSize: s.writeBufferSize,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}

		s.wsHandler.ServeWS(
			conn.UnderlyingConn(),
			newAbstractStream(conn),
			r.URL.Query(),
		)
	})
	s.httpHandler = mux

	go s.serve()

	return nil
}

// Stop will stop listening for new clients
func (s *Server) Stop() error {
	return s.listener.Close()
}
