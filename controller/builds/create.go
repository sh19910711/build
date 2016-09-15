package builds

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codestand/build/job"
	"github.com/codestand/build/jobqueue"
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

	log.Debug("Create: save source to tmp")
	r, _, err := c.Request.FormFile("file")
	if err != nil {
		respondError(c, err)
		return
	}
	if err := b.Save(r, "./tmp"); err != nil {
		respondError(c, err)
		return
	}

	// TODO: improve here
	b.Job.Callback = c.PostForm("callback")
	b.Job.Image = "build"
	b.Job.Commands = []string{"make"}
	b.Job.Artifacts = []string{"/app/app"}

	log.Debug("Create: push build job")
	go jobqueue.Push(b.Job)

	log.Debug("Create: success")
	job.Save(b.Job)
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
