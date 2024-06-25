package main

import (
	"gotunnels/tunnel"
	"log"
	"net"
	"net/http"
)

func waitForConnection() *tunnel.Tunnel {
	l, err := net.Listen("tcp", ":8090")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := l.Accept()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Tunnel established with host %s", conn.RemoteAddr())

	return tunnel.NewTunnel(conn)
}

func main() {
	tunnel := waitForConnection()

	http.Handle("/app/", tunnel)
	err := http.ListenAndServe(":9090", nil)
	log.Fatal(err)
}
