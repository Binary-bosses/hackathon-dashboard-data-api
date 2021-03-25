package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/Binary-bosses/hackathon-dashboard-data-api/server"
)

func main() {

	serverPath := "0.0.0.0:8080"
	log.Println("Starting server at " + serverPath)
	srv, err := server.NewServer()
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until a signal is received.
	go func() {
		<-c
		log.Println("Stopping server")
		srv.Stop()
	}()

	if err := srv.Start(serverPath); err != nil {
		log.Fatal(err)
	}
}
