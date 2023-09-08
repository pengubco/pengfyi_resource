package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	foo(3)
	foo(0)
}

func foo(i int) {
	defer fmt.Println("foo() in defer")

	switch i {
	case 0:
		os.Exit(0)
	case 1:
		os.Exit(1)
	case 2:
		log.Fatal("fatal")
	default:
		panic(i)
	}
}
