package controller

import (
	"github.com/gin-gonic/gin"
)

func Mount() {
	r := gin.Default()
	r.Run()
}
