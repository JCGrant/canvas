package main

import (
	"flag"

	"github.com/JCGrant/paint/server"
)

var address = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()
	s := server.New(*address)
	s.Serve()
}
