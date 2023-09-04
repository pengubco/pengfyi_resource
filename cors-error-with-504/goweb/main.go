package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	port := 8083

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		// sleep long enough to trigger gateway timeout.
		time.Sleep(3 * time.Second)

		log.Printf("%v \n%v\n", r.Method, r.Header)
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		fmt.Fprintf(w, "Hello, %v", time.Now())
	})

	url := fmt.Sprintf(":%d", port)
	log.Fatal(http.ListenAndServe(url, nil))
}
