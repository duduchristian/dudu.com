package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	start := time.Now()
	parallelPrint(10000)
	res1 := time.Since(start)

	start = time.Now()
	sequentialPrint(10000)
	res2 := time.Since(start)

	fmt.Println(res1)
	fmt.Println(res2)
}

func sequentialPrint(n int) {
	for i := 0; i < n; i++ {
		fmt.Println(i)
	}
}

func parallelPrint(n int) {
	c := make(chan int, 1)
	quit := make(chan struct{})
	var wg sync.WaitGroup

	c <- 0
	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
		loop:
			for {
				select {
				case num := <-c:
					if num < n {
						fmt.Println(num)
						c <- num + 1
					} else {
						close(quit)
					}
				case <-quit:
					break loop
				}
			}
		}()
	}

	wg.Wait()
}
