package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"os"
	"time"
)

func main() {
	timeout := flag.Int("timeout", 3, "timeout [seconds]")
	flag.Parse()

	r := gin.Default()

	go func() {
		time.Sleep(time.Duration(*timeout) * time.Second)
		os.Exit(1)
	}()

	r.POST("/callback", func(c *gin.Context) {
		os.Exit(0)
	})

	r.Run()
}
