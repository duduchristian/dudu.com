package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"sync"
)

var bufferPool *sync.Pool

const bufferKey = "BUFFER_TO_BE_PUT_BACK"

func init() {
	bufferPool = &sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, 512))
		},
	}
}

func ReadAll(context *gin.Context, r io.Reader) (b []byte, err error) {
	buffer := bufferPool.Get().(*bytes.Buffer)
	defer func() {
		if err != nil {
			buffer.Reset()
			bufferPool.Put(buffer)
			return
		}
		buffers, exists := context.Get(bufferKey)
		if exists {
			context.Set(bufferKey, append(buffers.([]*bytes.Buffer), buffer))
			return
		}
		context.Set(bufferKey, []*bytes.Buffer{buffer})
	}()

	_, err = buffer.ReadFrom(r)
	if err == nil {
		b = buffer.Bytes()
	}
	return
}

func WithBufferReleaser() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Next()
		buffers, exists := context.Get(bufferKey)
		if !exists {
			return
		}
		for _, buffer := range buffers.([]*bytes.Buffer) {
			buffer.Reset()
			bufferPool.Put(buffer)
		}
	}
}
