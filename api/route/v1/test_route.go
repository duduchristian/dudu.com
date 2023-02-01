package route

import (
	"fmt"
	"github.com/amitshekhariitbhu/go-backend-clean-architecture/bootstrap"
	"github.com/gin-gonic/gin"
	"time"
)

func NewTestRouter(env *bootstrap.Env, timeout time.Duration, group *gin.RouterGroup) {
	group.GET("/test", func(context *gin.Context) {
		fmt.Fprintf(context.Writer, "Hello World!")
	})
}
