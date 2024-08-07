package router

import (
	"log"
	"net/http"
)

func DefaultRouter() {
	http.HandleFunc("GET /moshe", func(w http.ResponseWriter, r *http.Request) {
		log.Println("sending Moshe to client")
		_, _ = w.Write([]byte("Moshe!\n"))
	})
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		log.Println("sending Hello World to client")
		_, _ = w.Write([]byte("Hello World!\n"))
	})
}
