package dudu_com

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/valyala/fasthttp"
	"log"
	"net/http"
	"sync"
	"syscall"
	"testing"
	"time"
)

func prepareServer(b *testing.B) int {
	pid, err := syscall.ForkExec("./cmd/main", []string{""}, nil)
	if err != nil {
		fmt.Println(err)
	}
	time.Sleep(300 * time.Millisecond)
	return pid
}

func prepareHertzServer(b *testing.B) int {
	pid, err := syscall.ForkExec("./cmd/main", []string{"", "-use-fasthttp"}, nil)
	if err != nil {
		fmt.Println(err)
	}
	time.Sleep(300 * time.Millisecond)
	return pid
}

func stopChildProcess(b *testing.B, pid int) {
	err := syscall.Kill(pid, syscall.SIGINT)
	if err != nil {
		b.Errorf("failed to kill child process: %v", err)
		return
	}
	wpid, err := syscall.Wait4(pid, nil, 0, nil)
	if err != nil {
		b.Errorf("failed to recycle child process: %v", err)
		return
	}
	if wpid != pid {
		b.Error("wpid is wrong")
		return
	}
	fmt.Printf("%d Done\n", pid)
}

func BenchmarkHttp(b *testing.B) {
	pid := prepareServer(b)

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

func BenchmarkFasthttp(b *testing.B) {
	pid := prepareServer(b)

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

var (
	respPool = sync.Pool{New: func() any {
		return &protocol.Response{}
	}}
)

func getResponse() *protocol.Response {
	return respPool.Get().(*protocol.Response)
}

func putResponse(res *protocol.Response) {
	res.Reset()
	respPool.Put(res)
}

func BenchmarkHertz(b *testing.B) {
	pid := prepareServer(b)

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
