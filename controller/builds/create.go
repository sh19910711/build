package builds

import (
	"github.com/codestand/build/model"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

// POST /builds
// - params[file] := tar-ball (required)
// - params[callback] := URL fired after completed build (required)
// - returns {"id": "<job-id>"}
func Create(c *gin.Context) {
	b := model.NewBuild()

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
	b.SaveJob()
	go b.PushJobQueue()

	c.JSON(http.StatusOK, gin.H{"id": b.Job.Id})
}

func respondError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{"msg": err})
}

func formFileReader(c *gin.Context) (io.Reader, error) {
	if r, _, err := c.Request.FormFile("file"); err != nil {
		return nil, err
	} else {
		return r, nil
	}
}
