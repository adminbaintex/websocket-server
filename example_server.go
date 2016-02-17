package main

import (
	"flag"
	"github.com/escrichov/websocket-server/server"
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

	log.Println(conn.RemoteAddr(), "CONNECTED", "Query params:", query)
	for {
		m, ok := <-s.Incoming()
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
	}
}

func main() {

	flag.Parse()

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	s := server.NewServer(&exampleHandler{})

	if err := s.ListenAndServe(*addr, *path); err != nil {
		log.Fatal(err)
	}

	log.Printf("Server listening on: %s, WS Path: %s", *addr, *path)

	// Terminate server in 60 seconds
	time.Sleep(time.Second * 60)
	if err := s.Stop(); err != nil {
		log.Fatal(err)
	}

	log.Println("Server stopped listening corretly")
}
