package dudu_com

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/valyala/fasthttp"
	"log"
	"net/http"
	"testing"
)

func BenchmarkHertzServer(b *testing.B) {
	pid := prepareHertzServer(b)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		c := &http.Client{}
		for pb.Next() {
			_, err := c.Get("http://localhost:8080/v1/test")
			if err != nil {
				b.Fatalf("Error: %v", err)
			}
		}
	})
	b.StopTimer()

	stopChildProcess(b, pid)
}

func BenchmarkHertzServerWithHertzClient(b *testing.B) {
	pid := prepareHertzServer(b)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		c, err := client.NewClient()
		if err != nil {
			b.Fatalf("Error:%v", err)
		}
		req := &protocol.Request{}
		req.SetMethod(consts.MethodGet)
		req.SetRequestURI("http://localhost:8080/v1/test")

		for pb.Next() {
			res := getResponse()
			err = c.Do(context.Background(), req, res)
			if err != nil {
				b.Fatalf("Error: %v", err)
			}
			putResponse(res)
		}
	})
	b.StopTimer()

	stopChildProcess(b, pid)
}

func BenchmarkHertzServerWithFasthttpClient(b *testing.B) {
	pid := prepareHertzServer(b)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		c := &fasthttp.HostClient{
			Addr: "localhost:8080",
		}
		for pb.Next() {
			statusCode, _, err := c.Get(nil, "http://localhost:8080/v1/test")
			if err != nil {
				log.Fatalf("Error when request through local proxy: %v", err)
			}
			if statusCode != fasthttp.StatusOK {
				log.Fatalf("Unexpected status code: %d. Expecting %d", statusCode, fasthttp.StatusOK)
			}
			if err != nil {
				b.Fatalf("Error: %v", err)
			}
		}
	})
	b.StopTimer()

	stopChildProcess(b, pid)
}
