package builds

import (
	"github.com/gin-gonic/gin"
)

func Mount(r *gin.Engine) {
	r.GET("/builds", Index)
	r.POST("/builds", Create)
	r.GET("/builds/:id", Show)
	r.GET("/builds/:id/log.txt", ShowLog) // TODO: has many
}
