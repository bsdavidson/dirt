package dirt

import (
	"fmt"
	"github.com/noonat/dirt/parser"
	"strings"
)

type Event interface {
	Process(s *Server) error
	String() string
}

type clientCommandEvent struct {
	client *Client
	input  string
	tokens []parser.Token
}

func (e *clientCommandEvent) Process(s *Server) error {
	if len(e.tokens) == 0 {
		e.client.writeChannel <- "> "
		return nil
	}
	cmdToken := e.tokens[0]
	cmd := strings.ToLower(cmdToken.Value)
	switch cmd {
	case "say":
		remainder := strings.TrimSpace(e.input[cmdToken.Start+cmdToken.Width:])
		e.client.writeChannel <- fmt.Sprintf("You say \"%s\"\n", remainder)
		for _, c := range s.Clients {
			if c != e.client {
				c.writeChannel <- fmt.Sprintf("Someone says \"%s\"\n", remainder)
			}
		}
	case "quit":
		e.client.writeChannel <- "Goodbye.\n"
		e.client.Close()
		return nil
	default:
		e.client.writeChannel <- "Huh? I don't understand...\n"
	}
	e.client.writeChannel <- "\n> "
	return nil
}

func (e *clientCommandEvent) String() string {
	return "clientCommandEvent"
}
