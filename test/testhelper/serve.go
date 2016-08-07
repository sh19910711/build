package testhelper

import (
	"github.com/gin-gonic/gin"
	"net/http/httptest"
)

func Serve(path string, handler func(c *gin.Context)) *httptest.Server {
	r := gin.Default()
	r.POST(path, handler)
	return httptest.NewServer(r)
}
