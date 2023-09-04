package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"time"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(2)

	http1Router := gin.Default()
	setRouter(http1Router)
	go func() {
		defer wg.Done()
		http1Router.Run(":8080")
	}()

	http2Router := gin.Default()
	setRouter(http2Router)
	go func() {
		defer wg.Done()
		http2Router.RunTLS(":8443", "./resources/server.pem", "./resources/server.key")
	}()
	wg.Wait()
}

func setRouter(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"content-type"},
	}))

	r.GET("/ping", func(c *gin.Context) {
		cnt := 5
		for {
			select {
			case <-c.Request.Context().Done():
				fmt.Println("context cancelled")
				return
			default:
				if cnt == 0 {
					c.String(http.StatusOK, "pong")
					return
				}
				cnt--
				fmt.Printf("processing.... ETA %d\n", cnt)
				time.Sleep(time.Second)
			}
		}
	})

	r.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "world")
	})
}
