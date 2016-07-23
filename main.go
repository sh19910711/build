package main

import log "github.com/Sirupsen/logrus"
import "github.com/gin-gonic/gin"

func main() {
	log.Info("starting build server")

	r := gin.Default()

	r.Run()
}
