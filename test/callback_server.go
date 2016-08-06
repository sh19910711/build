package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"os"
	"time"
)

func main() {
	// "--timeout 0" waits until request has been received
	timeout := flag.Int("timeout", 3, "timeout [seconds]")
	flag.Parse()

	r := gin.Default()
	exitCode := make(chan int)

	// timeout
	go func() {
		time.Sleep(time.Duration(*timeout) * time.Second)
		if *timeout > 0 {
			os.Exit(1)
		}
	}()

	// callback action
	r.POST("/callback", func(c *gin.Context) {
		c.JSON(200, gin.H{"hello": "world"})
		exitCode <- 0
	})

	// wait exitCode
	go func() {
		code := <-exitCode
		os.Exit(code)
	}()

	r.Run()
}
