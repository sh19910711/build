package main

import (
	"github.com/codestand/build/controller"
	"github.com/codestand/build/env"
	"github.com/codestand/build/jobqueue"
	"github.com/gin-gonic/gin"
)

func main() {
	jobqueue.Init()
	env.Init()

	go jobqueue.Wait()

	r := gin.Default()
	r.POST("/builds", controller.Create)
	r.Run()
}
