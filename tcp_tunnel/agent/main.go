package main

import (
	"gotunnels/router"
	"log"
	"net"
	"net/http"

	"github.com/hashicorp/yamux"
)

func Listen(managerAddr string) (net.Listener, error) {
	conn, err := net.Dial("tcp", managerAddr)
	if err != nil {
		return nil, err
	}

	return yamux.Client(conn, nil)
}

func main() {
	l, err := Listen("127.0.0.1:8090")
	if err != nil {
		log.Fatal(err)
	}

	router.DefaultRouter()

	err = http.Serve(l, nil)
	log.Fatal(err)
}
