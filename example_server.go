package main

import (
	"flag"
	"git.baintex.com/sentio/websocket-server/server"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"net/url"
	"time"
)

var addr = flag.String("addr", "127.0.0.1:8082", "ip:port to listening to")
var path = flag.String("path", "/app", "Path of the websocket application")

type exampleHandler struct{}

func (handler *exampleHandler) ServeWS(conn net.Conn, s server.Stream, query url.Values) {
	defer func() {
		log.Println(conn.RemoteAddr(), "CLOSED")
		s.Close()
	}()

	// New timer to close connection
	timer := time.NewTimer(time.Second * 3)

	log.Println(conn.RemoteAddr(), "CONNECTED", "Query params:", query)
	for {
		select {
		case m, ok := <-s.Incoming():
			if !ok {
				return
			}

			log.Println(conn.RemoteAddr(), "RECEIVED", m)

			msgSend := server.NewWSTextMessage([]byte("New text message"))
			err := s.Send(msgSend)
			if err != nil {
				log.Print(conn.RemoteAddr(), "ERROR", err)
				return
			}
			log.Println(conn.RemoteAddr(), "SENDED", msgSend)

			if err := s.Send(m); err != nil {
				log.Print(conn.RemoteAddr(), "ERROR", err)
				return
			}
			log.Println(conn.RemoteAddr(), "SENDED", m)
		case <-timer.C:
			log.Println("Connection exceed 3 second. Closing")
			return
		}
	}
}

func main() {

	flag.Parse()

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	readBufferSize := 1024
	writeBufferSize := 1024

	s := server.NewServer(&exampleHandler{}, readBufferSize, writeBufferSize)

	if err := s.ListenAndServe(*addr, *path); err != nil {
		log.Fatal(err)
	}

	log.Printf("Server listening on: %s, WS Path: %s", *addr, *path)

	// Terminate server in 30 seconds
	time.Sleep(time.Second * 30)
	if err := s.Stop(); err != nil {
		log.Fatal(err)
	}

	log.Println("Server stopped listening corretly")

	time.Sleep(time.Second * 10)
}
