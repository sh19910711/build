package builds

import (
	"github.com/codestand/build/model/build"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Index(c *gin.Context) {
	if all, err := build.All(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err})
	} else {
		c.JSON(http.StatusOK, gin.H{"builds": all})
	}
}
