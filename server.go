package dirt

import (
	"fmt"
	"log"
	"net"
)

type Server struct {
	Clients      []*Client
	listener     net.Listener
	closeChannel chan bool
	closed       bool
	events       chan Event
}

func NewServer() *Server {
	return &Server{
		closeChannel: make(chan bool, 1),
		events:       make(chan Event, 1024),
	}
}

func (s *Server) Close() {
	if !s.closed {
		s.closed = true
		s.closeChannel <- true
		close(s.closeChannel)
		s.listener.Close()
		for _, client := range s.Clients {
			client.Close()
		}
	}
}

func (s *Server) Emit(e Event) {
	s.events <- e
}

func (s *Server) ProcessEvents() {
	for {
		select {
		case e := <-s.events:
			if err := e.Process(s); err != nil {
				log.Printf("Error processing event %s: %s", e, err)
			}
		}
	}
}

func (s *Server) Run(listenAddr string) error {
	var err error
	s.listener, err = net.Listen("tcp", listenAddr)
	if err != nil {
		return fmt.Errorf("Error listening: %s", err.Error())
	}

	go s.ProcessEvents()

	log.Println("Listening on", listenAddr)
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return fmt.Errorf("Error accepting: %s", err.Error())
		}
		log.Println("Client connected")
		c := NewClient(conn, s)
		s.Clients = append(s.Clients, c)
		go c.Run()
	}
}
