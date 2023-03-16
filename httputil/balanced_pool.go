package httputil

import (
	"math/rand"
	"runtime"
	"sync"
	"unsafe"
)

type BanlancedPool struct {
	pools []*sync.Pool
	size  int
}

func NewBanlancedPool(new func() any) *BanlancedPool {
	size := runtime.NumCPU()
	size = 1
	bp := &BanlancedPool{
		pools: make([]*sync.Pool, size),
		size:  size,
	}
	for i := 0; i < size; i++ {
		bp.pools[i] = &sync.Pool{
			New: new,
		}
	}
	return bp
}

func (bp *BanlancedPool) Put(o interface{}) {
	index := int(uintptr(unsafe.Pointer(&o)) % uintptr(bp.size))
	bp.pools[index].Put(o)
}

func (bp *BanlancedPool) Get() interface{} {
	return bp.pools[rand.Intn(bp.size)].Get()
}
