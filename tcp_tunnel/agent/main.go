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
	ip := getIpFromCli()
	l, err := Listen(ip.String() + ":8090")
	if err != nil {
		log.Fatal(err)
	}

	router.DefaultRouter()

	err = http.Serve(l, nil)
	log.Fatal(err)
}

func getIpFromCli() net.IP {
	if len(os.Args) != 2 {
		fmt.Println("Usage: agent <IP_ADDRESS>")
		os.Exit(1)
	}
	ip := net.ParseIP(os.Args[1])
	if ip == nil {
		log.Fatal("Invalid IP address: ", ip)
	}

	return ip
}
