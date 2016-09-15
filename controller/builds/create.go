package builds

import (
	"github.com/codestand/build/model/build"
	"github.com/gin-gonic/gin"
	"net/http"
)

// POST /builds
// - params[file] := tar-ball (required)
// - params[callback] := URL fired after completed build (required)
// - returns {"id": "<job-id>"}
func Create(c *gin.Context) {
	b := build.New()

	r, _, err := c.Request.FormFile("file")
	if err != nil {
		respondError(c, err)
		return
	}
	if err := b.SaveSourceCode(r, "./tmp"); err != nil {
		respondError(c, err)
		return
	}

	b.SetCallbackURL(c.PostForm("callback"))
	b.SetWorker()
	build.Save(b)
	go b.PushJobQueue()

	c.JSON(http.StatusOK, b)
}
