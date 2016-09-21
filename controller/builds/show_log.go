package builds

import (
	"github.com/codestand/build/model/build"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

func ShowLog(c *gin.Context) {
	id := c.Param("id")
	if b, err := build.Find(id); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	} else {
		if logText, err := ioutil.ReadFile(b.Job.LogPath); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		} else {
			c.String(http.StatusOK, string(logText))
		}
	}
}
