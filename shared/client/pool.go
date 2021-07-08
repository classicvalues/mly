package client

import (
	"context"
	"crypto/tls"
	"golang.org/x/net/http2"
	"io"
	"net"
	"net/http"
	"reflect"
	"unsafe"
)

var emptyBuffer = make([][]byte, 0)
//This class is to workaround memory leak

type connPool struct {
	native http2.ClientConnPool
}

func (p *connPool) GetClientConn(req *http.Request, addr string) (*http2.ClientConn, error) {
	conn, err := p.native.GetClientConn(req, addr)
	return conn, err
}

func (p *connPool) MarkDead(conn *http2.ClientConn) {
	p.native.MarkDead(conn)
	//release free buf values
	connValue := reflect.ValueOf(conn).Elem()
	freeBufField := connValue.FieldByName("freeBuf")
	freeBufValue := reflect.NewAt(freeBufField.Type(), unsafe.Pointer(freeBufField.UnsafeAddr())).Elem().Interface()
	bs := freeBufValue.(*[][]byte)
	*bs = emptyBuffer
}



func newPool(host *Host, model string) *connPool {
	transport := &http2.Transport{
		AllowHTTP: true,
		DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
			return net.Dial(network, addr)
		},
	}
	httpClient := http.Client{
		Transport: transport,
		Timeout:   requestTimeout,
	}
	ctx, cancel := context.WithCancel(context.Background())
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, host.metaConfigURL(model), nil)
	if res, err := httpClient.Do(req); err != nil {
		if res.Body != nil {
			io.Copy(io.Discard, res.Body)
			res.Body.Close()
		}
	}
	cancel()
	httpClient.CloseIdleConnections()

	transportType := reflect.ValueOf(transport).Elem()
	poolField := transportType.FieldByName("connPoolOrDef")
	poolValue := reflect.NewAt(poolField.Type(), unsafe.Pointer(poolField.UnsafeAddr())).Elem().Interface()
	pool, _ := poolValue.(http2.ClientConnPool)
	return &connPool{native: pool}
}
