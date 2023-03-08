package main

import (
	"fmt"
	duducom "github.com/amitshekhariitbhu/go-backend-clean-architecture"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/domain"
	"github.com/valyala/fasthttp"
	"math/rand"
	"strings"
	"sync"
	"time"
)

var (
	logins = []*domain.LoginRequest{
		{
			Email:    "79560151@qq.com",
			Password: "123456",
		},
		{
			Email:    "795601511@qq.com",
			Password: "123456",
		},
		{
			Email:    "795601512@qq.com",
			Password: "123456",
		},
		{
			Email:    "795601513@qq.com",
			Password: "123456",
		},
	}
)

func main() {
	numWorker := 16
	times := make([]time.Duration, 0, numWorker)
	var wg sync.WaitGroup
	var timesLock sync.Mutex

	wg.Add(numWorker)
	for i := 0; i < numWorker; i++ {
		go func() {
			start := time.Now()
			doTest(1000)
			timesLock.Lock()
			times = append(times, time.Since(start))
			timesLock.Unlock()
			wg.Done()
		}()
	}
	wg.Wait()

	var total time.Duration
	for _, t := range times {
		total += t
	}
	fmt.Printf("average times is: %.2fs\n", total.Seconds()/float64(numWorker))
}

func doTest(count int) {
	c := &fasthttp.HostClient{
		Addr: "localhost:8080",
	}
	rand.Seed(time.Now().UnixNano())
	res := duducom.DoLogin(c, logins[rand.Intn(len(logins))])
	key := generateKey(res)

	for count > 0 {
		switch rand.Intn(3) {
		case 0:
			duducom.DoTest(c)
		case 1:
			duducom.DoProfile(c, key)
		case 2:
			tasks := duducom.DoGetTask(c, key)
			if len(tasks) < count {
				duducom.DoPostTask(c, key, "a task")
			}
		}
		count--
	}
}

func generateKey(resp *domain.LoginResponse) string {
	return strings.Join([]string{resp.RefreshToken, resp.AccessToken}, " ")
}
