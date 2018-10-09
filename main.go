package main

import (
	"flag"
	"log"

	"github.com/JCGrant/paint/server"
)

var address = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()
	s := server.New(*address)
	err := s.Serve()
	log.Fatalf("serving server failed: %s", err)
}
