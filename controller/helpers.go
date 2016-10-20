package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func atoi(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 32)
}

func internalError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{})
}

func notFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{})
}

func jsonResponse(c *gin.Context, h gin.H) {
	c.JSON(http.StatusOK, h)
}
