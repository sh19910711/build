package controller

import (
	"github.com/gin-gonic/gin"
)

func Mount() {
	r := gin.Default()
	r.POST("/builds/:id", Run)
	r.Run()
}
