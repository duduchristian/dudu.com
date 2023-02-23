package httputil

import (
	"testing"
	"unsafe"
)

func b2sUnsafe(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// go test -bench='^BenchmarkB2S' .
func BenchmarkB2S(b *testing.B) {
	url := []byte("1034275721057027507")
	for i := 0; i < b.N; i++ {
		b2s(url)
	}
}

func BenchmarkB2SUnsafe(b *testing.B) {
	url := []byte("1034275721057027507")
	for i := 0; i < b.N; i++ {
		b2sUnsafe(url)
	}
}
