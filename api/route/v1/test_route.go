package route

import (
	"bytes"
	"context"
	"fmt"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/bootstrap"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/httputil"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/gin-gonic/gin"
	"io"
	"time"
	"unsafe"
)

var bp = httputil.NewBanlancedPool(func() any {
	return bytes.NewBuffer(make([]byte, 0, 512))
})

func ReadAll(r io.Reader) (*bytes.Buffer, error) {
	b := bp.Get().(*bytes.Buffer)
	_, err := b.ReadFrom(r)
	return b, err
}

func s2b(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&struct {
		string
		Len int
		Cap int
	}{
		s,
		len(s),
		len(s),
	}))
}

func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func NewTestRouter(env *bootstrap.Env, timeout time.Duration, group *gin.RouterGroup) {
	group.GET("/test/:publisherId", func(context *gin.Context) {
		b, _ := ReadAll(context.Request.Body)
		defer func() {
			b.Reset()
			bp.Put(b)
		}()

		context.Request.Body.Close()
		fmt.Println(context.Request.URL.Query())
		fmt.Fprintf(context.Writer, b2s(b.Bytes()))
	})
	group.GET("/test1/:publisherId", func(context *gin.Context) {
		b, _ := io.ReadAll(context.Request.Body)

		context.Request.Body.Close()
		fmt.Println(context.Request.URL.Query())
		fmt.Fprintf(context.Writer, b2s(b))
	})
}

func NewTestRouterHertz(env *bootstrap.Env, timeout time.Duration, group *route.RouterGroup) {
	group.GET("/test", func(c context.Context, ctx *app.RequestContext) {
		fmt.Fprintf(ctx, "Hello World!")
	})
}
