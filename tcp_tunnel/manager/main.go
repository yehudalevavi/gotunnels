package main

import (
	"gotunnels/reverseproxy"
	"log"
	"net"
	"net/http"
)

func waitForTunnel() *reverseproxy.ReverseProxy {
	l, err := net.Listen("tcp", "127.0.0.1:8090")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := l.Accept()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Tunnel established with host %s", conn.RemoteAddr())

	return reverseproxy.NewReverseProxy(conn)
}

func main() {
	rp := waitForTunnel()

	http.Handle("/app/", rp)
	err := http.ListenAndServe(":9090", nil)
	log.Fatal(err)
}
