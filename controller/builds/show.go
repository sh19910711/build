package builds

import (
	"github.com/codestand/build/model/build"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GET /builds/<build-id>
// returns {id: <build-id>, finished: <>}
func Show(c *gin.Context) {
	id := c.Param("id")
	if b, err := build.Find(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err})
	} else {
		c.JSON(http.StatusOK, b)
	}
}
