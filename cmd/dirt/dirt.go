package main

import (
	"github.com/noonat/dirt"
	"log"
	"os"
	// "os/signal"
)

func main() {
	listenHost := os.Getenv("HOST")
	listenPort := os.Getenv("PORT")
	if listenPort == "" {
		listenPort = "4000"
	}
	listenAddr := listenHost + ":" + listenPort
	server := dirt.NewServer()
	err := server.Run(listenAddr)
	server.Close()
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
}
