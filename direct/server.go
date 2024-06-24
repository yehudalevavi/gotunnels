package main

import (
	"gotunnels/router"
	"log"
	"net"
	"net/http"
)

func main() {
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	router.DefaultRouter()
	err = http.Serve(l, nil)
	log.Fatal(err)
}
