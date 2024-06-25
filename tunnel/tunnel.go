package tunnel

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

type Tunnel struct {
	session *yamux.Session
	m       sync.Mutex
}

func NewTunnel(conn io.ReadWriteCloser) *Tunnel {
	rp := &Tunnel{}

	rp.NewConn(conn)

	return rp
}

func (t *Tunnel) NewConn(conn io.ReadWriteCloser) {
	session, err := yamux.Server(conn, nil)
	if err != nil {
		log.Fatal(err)
	}

	t.m.Lock()
	t.session = session
	t.m.Unlock()
}

func (t *Tunnel) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.m.Lock()
	session := t.session
	t.m.Unlock()
	if session == nil {
		http.Error(w, "Tunnel is not up", http.StatusInternalServerError)
		return
	} // empty session means no tunnel established

	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return session.Open()
		},
	}
	httpRP := &httputil.ReverseProxy{
		Transport: transport,
		Rewrite: func(pr *httputil.ProxyRequest) {
			target := &url.URL{Scheme: "http", Host: "manager RP", Path: "/"}
			pr.SetURL(target)
			pr.Out.URL.Path = strings.TrimPrefix(pr.In.URL.Path, "/agent")
		},
	}
	httpRP.ServeHTTP(w, r)
}
