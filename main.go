package main

import (
	"github.com/codestand/build/controller"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/", controller.Hello)
	r.POST("/builds", controller.RunBuild)

	r.Run()
}
