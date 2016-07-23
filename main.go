package main

import log "github.com/Sirupsen/logrus"
import "github.com/gin-gonic/gin"
import "net/http"

func main() {
	log.Info("starting build server")

	r := gin.Default()

	r.GET("/", func (c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"msg": "hello"})
	})

	r.Run()
}
