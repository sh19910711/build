package builds

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func respondError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{"msg": err})
}
