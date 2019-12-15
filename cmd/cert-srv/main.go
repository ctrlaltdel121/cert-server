package main

import (
	"log"

	"github.com/ctrlaltdel121/cert-server/srv"
)

func main() {
	s := srv.NewServer("file")
	err := s.Serve() // blocks
	if err != nil {
		log.Fatal(err)
	}
}
