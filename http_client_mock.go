package dudu_com

import (
	"encoding/json"
	"fmt"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/domain"
	"github.com/valyala/fasthttp"
	"log"
	"math/rand"
	"net/http"
)

const (
	urlPrefix = "http://localhost:8080/v1"
)

func DoLogin(c *fasthttp.HostClient, req *domain.LoginRequest) *domain.LoginResponse {
	var a fasthttp.Args
	a.Parse(fmt.Sprintf("email=%s&password=%s", req.Email, req.Password))

	statusCode, body, err := c.Post(nil, fmt.Sprintf("%s/login", urlPrefix), &a)
	if err != nil {
		panic(err)
	}
	if statusCode != http.StatusOK {
		panic("not ok")
	}

	var ret domain.LoginResponse
	err = json.Unmarshal(body, &ret)
	if err != nil {
		panic("Unmarshal error:" + err.Error())
	}
	return &ret
}

func DoProfile(c *fasthttp.HostClient, key string) *domain.Profile {
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()
	req.Header.Set("Authorization", key)
	req.SetRequestURI(fmt.Sprintf("%s/profile", urlPrefix))

	err := c.Do(req, res)
	if err != nil {
		panic(err)
	}
	statusCode := res.StatusCode()
	if statusCode != http.StatusOK {
		panic("not ok")
	}

	body := res.Body()
	var ret domain.Profile
	err = json.Unmarshal(body, &ret)
	if err != nil {
		panic("Unmarshal error:" + err.Error())
	}
	return &ret
}

func DoGetTask(c *fasthttp.HostClient, key string) []*domain.Task {
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(res)
	}()
	req.Header.Set("Authorization", key)
	req.SetRequestURI(fmt.Sprintf("%s/task", urlPrefix))

	err := c.Do(req, res)
	if err != nil {
		panic(err)
	}
	statusCode := res.StatusCode()
	if statusCode != http.StatusOK {
		panic("not ok")
	}

	body := res.Body()
	var ret []*domain.Task
	err = json.Unmarshal(body, &ret)
	if err != nil {
		panic("Unmarshal error:" + err.Error())
	}
	return ret
}

func DoPostTask(c *fasthttp.HostClient, key, title string) {
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(res)
	}()
	req.Header.SetMethod("POST")
	req.Header.Set("Authorization", key)
	req.Header.Set("Content-Type", "application/json")
	req.SetBodyString(fmt.Sprintf("{\"title\":\"%s\"}", title))
	req.SetRequestURI(fmt.Sprintf("%s/task", urlPrefix))

	err := c.Do(req, res)
	if err != nil {
		panic(err)
	}
	statusCode := res.StatusCode()
	if statusCode != http.StatusOK {
		panic(string(res.Body()))
	}
}

func generateRandomBytes() []byte {
	length := rand.Intn(4096)
	b := make([]byte, length)
	for i := 0; i < length; i++ {
		b[i] = 'a'
	}
	return b
}

func DoTest(c *fasthttp.HostClient) {
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(res)
	}()

	req.SetRequestURI(fmt.Sprintf("%s/test/pub10000", urlPrefix))
	req.Header.SetMethod("GET")
	req.Header.Set("Content-Type", "application/json")
	if rand.Intn(2) == 0 {
		req.Header.Set("Test-Flag", "dudu")
	} else {
		req.Header.Set("Test-Flag", "christian")
	}
	req.SetBody(generateRandomBytes())
	err := c.Do(req, res)
	statusCode := res.StatusCode()
	if err != nil {
		log.Fatalf("Error when request through local proxy: %v", err)
	}
	if statusCode != fasthttp.StatusOK {
		log.Fatalf("Unexpected status code: %d. Expecting %d", statusCode, fasthttp.StatusOK)
	}
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
