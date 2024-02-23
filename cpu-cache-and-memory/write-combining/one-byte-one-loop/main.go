package main

import (
	"fmt"
	"time"
)

type A [1]uint8

func main() {
	arraySize := 1 << 20
	d0 := make([]A, arraySize)
	d1 := make([]A, arraySize)
	d2 := make([]A, arraySize)
	d3 := make([]A, arraySize)
	d4 := make([]A, arraySize)
	d5 := make([]A, arraySize)
	d6 := make([]A, arraySize)
	d7 := make([]A, arraySize)

	// writes a constant, 2.
	var value A = [1]uint8{2}
	loopCount := 1_000_000_000

	mask := arraySize - 1
	t := time.Now()
	for i := 0; i < loopCount; i++ {
		j := i & mask
		d0[j] = value
		d1[j] = value
		d2[j] = value
		d3[j] = value
		d4[j] = value
		d5[j] = value
		d6[j] = value
		d7[j] = value
	}
	fmt.Printf("takes %d ns\n", time.Since(t).Nanoseconds())
}
