package main

import (
	"github.com/codestand/build/controller/builds"
	_ "github.com/codestand/build/env"
	"github.com/codestand/build/jobqueue"
	"github.com/gin-gonic/gin"
)

func main() {
	go jobqueue.Wait()
	defer jobqueue.Close()

	r := gin.Default()
	builds.Mount(r)
	r.Run()
}
