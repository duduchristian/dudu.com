package main

import (
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/httputil"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/metrics"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/prefork"
	"os"
	"time"

	routeV1 "github.com/amitshekhariitbhu/go-backend-clean-architecture/api/route/v1"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/bootstrap"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

const (
	useHertzFlag    = "-use-hertz"
	useFastHttpFlag = "-use-fasthttp"
	usePreforkFlag  = "-use-prefork"
)

func main() {
	app := bootstrap.App()

	env := app.Env

	db := app.Mongo.Database(env.DBName)
	defer app.CloseDBConnection()

	timeout := time.Duration(env.ContextTimeout) * time.Second

	if !useHertzHttpServer() {
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
			fs := httputil.NewFasthttpServer(httputil.NewFastHTTPHandler(gin.Handler()))
			if err := fs.ListenAndServe(addr); err != nil {
				panic(err)
			}
			if err := fasthttp.ListenAndServe(addr, httputil.NewFastHTTPHandler(gin.Handler())); err != nil {
				panic(err)
			}
		}
	} else {
		h := server.Default(config.Option{F: func(o *config.Options) {
			o.Addr = env.ServerAddress
		}})
		v1 := h.Group("/v1")
		routeV1.NewTestRouterHertz(env, timeout, v1)
		h.Spin()
	}
}

func useHertzHttpServer() bool {
	for _, arg := range os.Args {
		if arg == useHertzFlag {
			return true
		}
	}
	return false
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
