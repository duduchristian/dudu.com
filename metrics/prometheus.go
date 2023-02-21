package metrics

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"path"
	"sync"
)

var opsProcessedInBidder = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "processed_ops_total",
	Help: "The total number of processed events",
}, []string{"path"})

func InitRouter(router gin.IRouter) {
	prometheus.MustRegister(opsProcessedInBidder)
	router.Use(WithPublisherRpsMetric(opsProcessedInBidder))
	router.GET("/metrics", func(context *gin.Context) {
		promhttp.Handler().ServeHTTP(context.Writer, context.Request)
	})
}

var stringMap = map[string]struct{}{}
var lock sync.Mutex

func WithPublisherRpsMetric(ops *prometheus.CounterVec) gin.HandlerFunc {
	return func(context *gin.Context) {
		apiPath := context.Param("publisherId")
		if apiPath == "" {
			apiPath = path.Base(context.FullPath())
		}
		//metrics.GetOrCreateCounter(fmt.Sprintf(`bidder_processed_ops_total{path="%s"}`, apiPath)).Inc()
		counter := ops.WithLabelValues(apiPath)
		counter.Inc()
		lock.Lock()
		stringMap[context.Request.Header.Get("Test-Flag")] = struct{}{}
		if repeat, ok := checkoutRepeat(stringMap); ok {
			fmt.Println(repeat)
		}
		lock.Unlock()
	}
}

func checkoutRepeat(m map[string]struct{}) (string, bool) {
	s := map[string]struct{}{}
	for k := range m {
		if _, ok := s[k]; ok {
			return k, true
		}
		s[k] = struct{}{}
	}
	return "", false
}
