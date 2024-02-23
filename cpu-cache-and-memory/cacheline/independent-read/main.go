package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	t := time.Now()
	var a [2]A
	cnt := 1_000_000_000

	var wg sync.WaitGroup
	wg.Add(2)
	for i := 0; i < 2; i++ {
		go func(p *A) {
			defer wg.Done()
			sum := int64(0)
			for j := 0; j < cnt; j++ {
				sum += p.x
				sum++
			}
			fmt.Println(sum)
		}(&a[i])
	}
	wg.Wait()
	fmt.Printf("takes %d ns\n", time.Since(t).Nanoseconds())
}

type A struct {
	x   int64
	pad [7]int64
}
