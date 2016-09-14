package main

import (
	"github.com/codestand/build/controller"
	_ "github.com/codestand/build/env"
	"github.com/codestand/build/jobqueue"
	"github.com/gin-gonic/gin"
)

func main() {
	go jobqueue.Wait()

	r := gin.Default()
	r.POST("/builds", controller.Create)
	r.Run()
}
