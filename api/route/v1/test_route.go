package route

import (
	"context"
	"fmt"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/bootstrap"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/gin-gonic/gin"
	"time"
)

func NewTestRouter(env *bootstrap.Env, timeout time.Duration, group *gin.RouterGroup) {
	group.GET("/test", func(context *gin.Context) {
		fmt.Fprintf(context.Writer, "Hello World!")
	})
}

func NewTestRouterHertz(env *bootstrap.Env, timeout time.Duration, group *route.RouterGroup) {
	group.GET("/test", func(c context.Context, ctx *app.RequestContext) {
		fmt.Fprintf(ctx, "Hello World!")
	})
}
