package builds

import (
	"github.com/gin-gonic/gin"
)

func Mount(r *gin.Engine) {
	r.POST("/builds", Create)
	r.GET("/builds/:id", Show)
}
