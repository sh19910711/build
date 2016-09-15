// +build release

package env

const DEBUG = false

func init() {
	gin.SetMode(gin.ReleaseMode)
}
