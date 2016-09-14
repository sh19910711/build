package main

import (
	"github.com/codestand/build/controller"
	_ "github.com/codestand/build/env"
	"github.com/codestand/build/jobqueue"
	"github.com/gin-gonic/gin"
)

func main() {
	go jobqueue.Wait()
	defer jobqueue.Close()

	r := gin.Default()
	controller.MountBuilds(r)
	r.Run()
}
