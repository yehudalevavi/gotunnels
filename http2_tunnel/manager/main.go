package main

import (
	"context"
	"fmt"
	"gotunnels/reverseproxy"
	"log"
	"net/http"
)

type conn struct {
	w  http.ResponseWriter
	r  *http.Request
	rc *http.ResponseController
}

func newConn(w http.ResponseWriter, r *http.Request) *conn {
	return &conn{w: w, r: r, rc: http.NewResponseController(w)}
}

func (c *conn) Read(b []byte) (n int, err error) {
	return c.r.Body.Read(b)
}

func (c *conn) Write(b []byte) (n int, err error) {
	n, err = c.w.Write(b)
	if err == nil {
		err = c.rc.Flush()
	}

	return n, err
}

func (c *conn) Close() error {
	return c.r.Body.Close()
}

func main() {
	rp := &reverseproxy.ReverseProxy{}
	ctx := context.Background()

	http.HandleFunc("/connect", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/connect called")
		if r.ProtoMajor != 2 {
			http.Error(w, "insufficient http version", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		rp.NewConn(newConn(w, r))
		<-ctx.Done()
		fmt.Println("Context canceled, tunnel is down")
	})
	http.Handle("/app/", rp)

	go func() {
		err := http.ListenAndServeTLS(":9090", "./testdata/cert.pem", "./testdata/key.pem", nil)
		log.Fatal(err)
	}()

	<-ctx.Done()
}
