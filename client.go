// FIXME: remove client from server after disconnection
// FIXME: don't close channels until client is removed from server

package dirt

import (
	"bufio"
	"github.com/noonat/dirt/parser"
	"io"
	"log"
	"net"
)

// Client represents a single user connected to the server.
type Client struct {
	*Entity
	net.Conn
	*Server
	closeChannel chan bool
	readChannel  chan string
	writeChannel chan string
	closed       bool
}

// NewClient creates a new client object for the given connection.
func NewClient(conn net.Conn, server *Server) *Client {
	return &Client{
		Conn:         conn,
		Server:       server,
		closeChannel: make(chan bool, 3),
		readChannel:  make(chan string),
		writeChannel: make(chan string),
	}
}

// Close the connection.
func (c *Client) Close() {
	if !c.closed {
		log.Println("Client disconnected")
		c.closed = true
		for i := 0; i < 3; i++ {
			c.closeChannel <- true
		}
		close(c.closeChannel)
		close(c.readChannel)
		close(c.writeChannel)
		c.Conn.Close()
	}
}

func (c *Client) Run() {
	go c.parseMessages()
	go c.readMessages()
	go c.writeMessages()
}

func (c *Client) parseMessages() {
	for {
		select {
		case <-c.closeChannel:
			return
		case msg := <-c.readChannel:
			tokens, err := parser.Parse(msg)
			if err != nil {
				log.Printf("Error parsing message from client: %s (message was %+v)\n", err.Error(), msg)
				break
			}
			c.Server.Emit(&clientCommandEvent{client: c, input: msg, tokens: tokens})
		}
	}
}

func (c *Client) readMessages() {
	r := bufio.NewReader(c.Conn)
	for {
		select {
		case <-c.closeChannel:
			return
		default:
		}
		s, err := r.ReadString('\n')
		if err == io.EOF {
			c.Close()
			return
		} else if err != nil {
			if !c.closed {
				log.Println("Error reading from client:", err.Error())
			}
			c.Close()
			return
		}
		c.readChannel <- s
	}
}

func (c *Client) writeMessages() {
	for {
		select {
		case <-c.closeChannel:
			return
		case msg := <-c.writeChannel:
			c.Conn.Write([]byte(msg))
		}
	}
}
