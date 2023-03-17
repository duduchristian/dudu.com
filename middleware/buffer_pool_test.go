package middleware

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
	"unsafe"
)

func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func startServer() *httptest.Server {
	engine := gin.Default()
	engine.POST("/test", func(context *gin.Context) {
		b, err := ioutil.ReadAll(context.Request.Body)
		context.Request.Body.Close()
		if err != nil {
			_ = context.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		fmt.Fprintf(context.Writer, b2s(b))
	})
	return httptest.NewServer(engine)
}

func startServerWithBufferPool() *httptest.Server {
	engine := gin.Default()
	engine.Use(WithBufferReleaser())
	engine.POST("/test", func(context *gin.Context) {
		b, err := ReadAll(context, context.Request.Body)
		context.Request.Body.Close()
		if err != nil {
			_ = context.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		fmt.Fprintf(context.Writer, b2s(b))
	})
	return httptest.NewServer(engine)
}

var pool = &sync.Pool{
	New: func() interface{} { return make([]byte, 0, 4096) },
}

func generateRandomBytes() []byte {
	length := rand.Intn(4096)
	b := pool.Get().([]byte)
	b = b[:length]
	return b
}

// go test -run "^TestReadAll$" .
func TestReadAll(t *testing.T) {
	gin.DefaultWriter = io.Discard
	server := startServerWithBufferPool()
	defer server.Close()
	uri := fmt.Sprintf("%s/test", server.URL)

	for i := 0; i < 1000; i++ {
		data := generateRandomBytes()
		for j := 0; j < len(data); j++ {
			data[j] = 'a' + byte(rand.Intn(26))
		}
		resp, err := http.Post(uri, "text/plain", bytes.NewBuffer(data))
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Error response with status code: %d", resp.StatusCode)
		}
		resBytes, _ := ioutil.ReadAll(resp.Body)
		if b2s(resBytes) != b2s(data) {
			t.Fatalf("Incorrect response data")
		}
		resp.Body.Close()
		pool.Put(data)
	}
}

func benchmarkRequestServer(b *testing.B, server *httptest.Server) {
	uri := fmt.Sprintf("%s/test", server.URL)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data := generateRandomBytes()
		resp, err := http.Post(uri, "text/plain", bytes.NewBuffer(data))
		if err != nil {
			b.Fatalf("Error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			b.Fatalf("Error response with status code: %d", resp.StatusCode)
		}
		resp.Body.Close()
		pool.Put(data)
	}
	b.StopTimer()
}

// go test -bench="^BenchmarkIoutilReadAll$" -benchmem .
// BenchmarkIoutilReadAll-8            9070            172326 ns/op           25214 B/op        154 allocs/op
func BenchmarkIoutilReadAll(b *testing.B) {
	gin.DefaultWriter = io.Discard
	server := startServer()
	defer server.Close()
	time.Sleep(time.Second)
	benchmarkRequestServer(b, server)
}

// go test -bench="^BenchmarkReadAllWithBufferPool$" -benchmem .
// BenchmarkReadAllWithBufferPool-8            9650            175071 ns/op           19295 B/op        153 allocs/op
func BenchmarkReadAllWithBufferPool(b *testing.B) {
	gin.DefaultWriter = io.Discard
	server := startServerWithBufferPool()
	defer server.Close()
	time.Sleep(time.Second)
	benchmarkRequestServer(b, server)
}
