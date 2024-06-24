package main

import (
	"crypto/tls"
	"errors"
	"gotunnels/router"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"

	"github.com/hashicorp/yamux"
	"golang.org/x/net/http2"
)

type conn struct {
	w *io.PipeWriter
	r io.ReadCloser
}

func (c *conn) Write(b []byte) (n int, err error) { return c.w.Write(b) }
func (c *conn) Read(b []byte) (n int, err error)  { return c.r.Read(b) }
func (c *conn) Close() error                      { return errors.Join(c.w.Close(), c.r.Close()) }

func TunnelDial(managerAddr string) (io.ReadWriteCloser, error) {
	cli := http.Client{
		Transport: &http2.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	pipeR, pipeW := io.Pipe()
	tunnelUrl := url.URL{Scheme: "https", Host: managerAddr, Path: "connect"}

	log.Println("Establishing tunnel connection with ", tunnelUrl.String())
	resp, err := cli.Post(tunnelUrl.String(), "", io.NopCloser(pipeR))
	if err != nil {
		return nil, err
	}

	return &conn{pipeW, resp.Body}, nil
}

func Listen(managerAddr string) (net.Listener, error) {
	tunnel, err := TunnelDial(managerAddr)
	if err != nil {
		return nil, err
	}

	return yamux.Client(tunnel, nil)
}

func main() {
	l, err := Listen("127.0.0.1:9090")
	if err != nil {
		log.Fatal(err)
	}

	router.DefaultRouter()

	err = http.Serve(l, nil)
	log.Fatal(err)
}