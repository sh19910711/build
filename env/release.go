// +build release

package env

const DEBUG = false

func Init() {
	gin.SetMode(gin.ReleaseMode)
}
