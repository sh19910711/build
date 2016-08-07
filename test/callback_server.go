package main

import (
	"archive/tar"
	"flag"
	"github.com/gin-gonic/gin"
	"io"
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
		r, _, err := c.Request.FormFile("file")
		if err != nil {
			panic(err)
		}

		// extract artifacts
		tr := tar.NewReader(r)
		for {
			header, err := tr.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				panic(err)
			}

			f, err := os.OpenFile("./tmp/"+header.Name, os.O_WRONLY|os.O_CREATE, os.FileMode(header.Mode))
			if err != nil {
				panic(err)
			}
			defer f.Close()

			if _, err := io.Copy(f, tr); err != nil {
				panic(err)
			}
		}

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
