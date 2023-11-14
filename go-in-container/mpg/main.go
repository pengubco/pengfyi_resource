package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

/*
This program is used to start many goroutines and wait for them to finish.
1. 1/3 goroutines sum integers, sleep a milisecond, and repeat.
2. 1/3 goroutines sleep for 1 minute.
3. 1/3 goroutines sum random integers.
*/
func main() {
	var numOfGoroutines int
	flag.IntVar(&numOfGoroutines, "numberOfGoroutines", 1, "number of goroutines to create")
	var quiet bool
	flag.BoolVar(&quiet, "quiet", false, "whether print out information")

	flag.Parse()

	if !quiet {
		fmt.Printf("go routine: %d\n", numOfGoroutines)
		fmt.Printf("NumCPU: %d, GOMAXPROCS: %d\n", runtime.NumCPU(), runtime.GOMAXPROCS(0))
	}

	var wg sync.WaitGroup
	m := numOfGoroutines / 3
	wg.Add(3 * m)
	for i := 0; i < m; i++ {
		go work1(&wg)
	}
	for i := 0; i < m; i++ {
		go work2(&wg)
	}
	for i := 0; i < m; i++ {
		go work3(&wg)
	}

	wg.Wait()
}

func work1(wg *sync.WaitGroup) {
	defer wg.Done()
	sum := 0
	for j := 0; j < 100000000; j++ {
		time.Sleep(time.Millisecond)
		sum += j
	}
}

func work2(wg *sync.WaitGroup) {
	defer wg.Done()
	<-time.After(time.Minute)
}

func work3(wg *sync.WaitGroup) {
	defer wg.Done()
	rand.Seed(0)
	result := 0
	for i := 0; i < 100000; i++ {
		result += rand.Intn(math.MaxInt)
	}
}
