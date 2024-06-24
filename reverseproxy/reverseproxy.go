package reverseproxy

import (
	"context"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"

	"github.com/hashicorp/yamux"
)

type ReverseProxy struct {
	session *yamux.Session
	m       sync.Mutex
}

func NewReverseProxy(conn io.ReadWriteCloser) *ReverseProxy {
	rp := &ReverseProxy{}

	rp.NewConn(conn)

	return rp
}

func (rp *ReverseProxy) NewConn(conn io.ReadWriteCloser) {
	session, err := yamux.Server(conn, nil)
	if err != nil {
		log.Fatal(err)
	}

	rp.m.Lock()
	rp.session = session
	rp.m.Unlock()
}

func (rp *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rp.m.Lock()
	session := rp.session
	rp.m.Unlock()

	if session == nil {
		http.Error(w, "Tunnel is not up", http.StatusInternalServerError)
		return
	}

	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return session.Open()
		},
	}
	httpRP := &httputil.ReverseProxy{
		Transport: transport,
		Rewrite: func(pr *httputil.ProxyRequest) {
			target := &url.URL{Scheme: "http", Host: "yamux", Path: "/"}
			pr.SetURL(target)
			pr.Out.URL.Path = strings.TrimPrefix(pr.In.URL.Path, "/app")
		},
	}
	httpRP.ServeHTTP(w, r)
}
