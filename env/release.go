// +build release

package env

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

const DEBUG = false

func init() {
	log.SetLevel(log.DebugLevel)
	log.Info("build server started in the release mode")
	gin.SetMode(gin.ReleaseMode)
}
