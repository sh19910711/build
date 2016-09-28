// +build release

package env

import (
	"github.com/gin-gonic/gin"
)

const DEBUG = false

func init() {
	gin.SetMode(gin.ReleaseMode)
}
