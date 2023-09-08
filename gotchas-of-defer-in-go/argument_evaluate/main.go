package main

import (
	"fmt"
	"log"
)

func main() {
	foo(1)

	fooFixed(1)
}

func foo(i int) (err error) {
	defer handleError(err)
	err = fmt.Errorf("error kind %d", i)
	return err
}

func handleError(err error) {
	log.Println(err)
}

func fooFixed(i int) (err error) {
	defer func() { // capture value using closure
		handleError(err)
	}()
	err = fmt.Errorf("error kind %d", i)
	return err
}
