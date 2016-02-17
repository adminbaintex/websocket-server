package server

import (
	"log"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
)

// WSHandler will receive new connections as streams.
type WSHandler interface {
	Serve(net.Conn, Stream)
}

// Server manages multiple Configurations and yields new connection as
// streams to the Handler.
type Server struct {
	// The Handler that receives new Streams.
	wsHandler WSHandler

	// Server Listener and Handler
	listener    net.Listener
	httpHandler http.Handler
}

// NewServer returns a new Server.
func NewServer(handler WSHandler) *Server {
	return &Server{
		wsHandler: handler,
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
	s.listener = l

	log.Printf("Server listening on: %s, WS Path: %s", address, path)

	mux := http.NewServeMux()
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}

		s.wsHandler.Serve(conn.UnderlyingConn(), newAbstractStream(conn))
	})
	s.httpHandler = mux

	go s.serve()

	return nil
}

// Stop will stop listening for new clients
func (s *Server) Stop() error {
	return s.listener.Close()
}
