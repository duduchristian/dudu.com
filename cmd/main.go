package main

import (
	"flag"
	"fmt"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/httputil"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/metrics"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/tuner"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/prefork"
	"io"
	"os"
	"time"

	routeV1 "github.com/amitshekhariitbhu/go-backend-clean-architecture/api/route/v1"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/bootstrap"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

const (
	useFastHttpFlag = "-use-fasthttp"
	usePreforkFlag  = "-use-prefork"
	useTunerFlag    = "-use-tuner"
)

func initProcess() {
	var (
		inCgroup = true
		percent  = 90.0
	)
	tuner.NewTuner(inCgroup, percent)
}

func main() {
	if useTuner() {
		fmt.Println("Use tuner!!!")
		initProcess()
	}
	var logDir string
	flag.StringVar(&logDir, "log_dir", "./log", "Specify log directory")

	app := bootstrap.App()

	env := app.Env

	db := app.Mongo.Database(env.DBName)
	defer app.CloseDBConnection()

	timeout := time.Duration(env.ContextTimeout) * time.Second

	f, err := os.OpenFile(fmt.Sprintf("%s/app01.log", logDir), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	gin.DefaultWriter = io.MultiWriter(f)

	gin := gin.Default()
	metrics.InitRouter(gin)
	routerV1 := gin.Group("v1")
	routeV1.Setup(env, timeout, db, routerV1)
	addr := env.ServerAddress
	pprof.Register(gin)
	if !useFastHttpServer() {
		gin.Run(addr)
	} else {
		if usePrefork() {
			s := &fasthttp.Server{
				Handler: httputil.NewFastHTTPHandler(gin.Handler()),
			}
			p := prefork.New(s)
			if err := p.ListenAndServe(addr); err != nil {
				panic(err)
			}
			return
		}
		if err := fasthttp.ListenAndServe(addr, httputil.NewFastHTTPHandler(gin.Handler())); err != nil {
			panic(err)
		}
		fs := httputil.NewFasthttpServer(httputil.NewFastHTTPHandler(gin.Handler()))
		if err := fs.ListenAndServe(addr); err != nil {
			panic(err)
		}
	}
}

func useFastHttpServer() bool {
	for _, arg := range os.Args {
		if arg == useFastHttpFlag {
			return true
		}
	}
	return false
}

func usePrefork() bool {
	for _, arg := range os.Args {
		if arg == usePreforkFlag {
			return true
		}
	}
	return false
}

func useTuner() bool {
	for _, arg := range os.Args {
		if arg == useTunerFlag {
			return true
		}
	}
	return false
}
