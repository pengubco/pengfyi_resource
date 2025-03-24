package main

import (
	"fmt"
	"time"
)

func main() {
	var d [][]byte
	time.Sleep(3 * time.Second) // sleep 3 seconds for the collection script to find PID.
	for range 200 {
		d = append(d, make([]byte, 8*50*1024*1024)) // 400 MiB
		time.Sleep(time.Second)
		// Reference d, in case compiler optimizes d away.
		fmt.Println(len(d))
	}
}
