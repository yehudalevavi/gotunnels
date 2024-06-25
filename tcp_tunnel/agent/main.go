package main

import (
	"fmt"
	"gotunnels/router"
	"log"
	"net"
	"net/http"
	"os"

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
	if len(os.Args) != 2 {
		fmt.Println("Usage: agent <IP_ADDRESS>")
		os.Exit(1)
	}
	ip := os.Args[1]
	if net.ParseIP(ip) == nil {
		log.Fatal("Invalid IP address: ", ip)
	}

	l, err := Listen(ip + ":8090")
	if err != nil {
		log.Fatal(err)
	}

	router.DefaultRouter()

	err = http.Serve(l, nil)
	log.Fatal(err)
}
