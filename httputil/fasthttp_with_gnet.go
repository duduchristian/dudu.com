package httputil

import (
	"github.com/panjf2000/gnet"
	"github.com/valyala/fasthttp"
	"math/rand"
	"net"
	"sync"
)

const NumWorker = 1

type FasthttpServer struct {
	*gnet.EventServer
	lns    [NumWorker]*InmemoryListener
	server *fasthttp.Server
	lock   sync.Mutex
	m      map[string]net.Conn
}

func NewFasthttpServer(handler fasthttp.RequestHandler) *FasthttpServer {
	s := &fasthttp.Server{}
	s.Handler = handler
	fs := &FasthttpServer{
		server: s,
	}
	for i := 0; i < NumWorker; i++ {
		ln := NewInmemoryListener()
		serverAddr := &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 12345}
		ln.SetLocalAddr(serverAddr)
		fs.lns[i] = ln
	}
	return fs
}

func (fs *FasthttpServer) ListenAndServe(addr string) error {
	for i := 0; i < NumWorker; i++ {
		go func(index int) {
			_ = fs.server.Serve(fs.lns[index])
		}(i)
	}
	return gnet.Serve(fs, addr, gnet.WithMulticore(true), gnet.WithReusePort(true))
}

var bytesPool = &sync.Pool{
	New: func() any {
		return make([]byte, 4096)
	},
}

func getBytes() []byte {
	return bytesPool.Get().([]byte)
}

func putBytes(b []byte) {
	bytesPool.Put(b)
}

func (fs *FasthttpServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	ln := fs.lns[rand.Intn(NumWorker)]
	conn, _ := ln.DialWithLocalAddr(c.RemoteAddr())
	defer conn.Close()
	conn.Write(frame)

	out, _ = conn.(*pipeConn).ReadAll()

	return
}
