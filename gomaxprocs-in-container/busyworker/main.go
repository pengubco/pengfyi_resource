package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	var numberOfBusyWorkers int
	var busyWorkSize int
	var numberOfOtherWorkers int
	var quiet bool
	flag.IntVar(&numberOfBusyWorkers, "numberOfBusyWorkers", 10, "number of busy workers. each worker is a goroutine")
	flag.IntVar(&busyWorkSize, "busyWorkSize", 1000000, "Number of random integer each worker sums")
	flag.IntVar(&numberOfOtherWorkers, "numberOfOtherWorkers", 1000, "number of other workers. each worker is a goroutine")
	flag.BoolVar(&quiet, "quiet", false, "whether print out information")
	flag.Parse()

	if !(quiet) {
		fmt.Printf("NumCPU: %d, GOMAXPROCS: %d\n", runtime.NumCPU(), runtime.GOMAXPROCS(0))
		fmt.Printf("number of busy workers: %d\n", numberOfBusyWorkers)
		fmt.Printf("number of busy work size: %d\n", busyWorkSize)
	}

	var wg sync.WaitGroup
	m := numberOfOtherWorkers / 2
	wg.Add(2 * m)
	for i := 0; i < m; i++ {
		go work1(&wg)
	}
	for i := 0; i < m; i++ {
		go work2(&wg)
	}

	var busyWorkerWG sync.WaitGroup
	busyWorkerWG.Add(numberOfBusyWorkers)

	var totalLatency atomic.Int64
	for i := 0; i < numberOfBusyWorkers; i++ {
		go func() {
			latency := busyWork(&busyWorkerWG, busyWorkSize)
			if !quiet {
				fmt.Printf("%.2f\n", latency.Seconds())
			}
			totalLatency.Add(latency.Milliseconds())
		}()
	}

	busyWorkerWG.Wait()
	fmt.Printf("average time taken %.1f\n",
		float64(totalLatency.Load()/int64(numberOfBusyWorkers))/float64(1000))
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
	<-time.After(time.Minute)
	<-time.After(time.Minute)
}

func busyWork(wg *sync.WaitGroup, busyWorkSize int) time.Duration {
	defer wg.Done()
	rand.Seed(0)
	result := 0
	t0 := time.Now()
	for i := 0; i < busyWorkSize; i++ {
		result += rand.Intn(math.MaxInt)
	}
	return time.Now().Sub(t0)
}
