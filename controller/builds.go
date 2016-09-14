package controller

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codestand/build/jobqueue"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type Build struct {
	job jobqueue.Job
}

func (b *Build) saveSource(c *gin.Context, prefix string) error {
	if err := os.MkdirAll(prefix, 0700); err != nil {
		return err
	}

	r, err := formFileReader(c)
	if err != nil {
		return err
	}
	src := filepath.Join(prefix, b.job.Id)
	if err := save(r, src); err != nil {
		return err
	} else {
		b.job.Src = src
	}

	return nil
}

func (b *Build) setup(c *gin.Context) {
	b.job.Callback = c.PostForm("callback")
	b.job.Image = "build"
	b.job.Commands = []string{"make"}
	b.job.Artifacts = []string{"/app/app"}
}

// POST /builds
// - params[file] := tar-ball (required)
// - params[callback] := URL fired after completed build (required)
// - returns {"id": "<job-id>"}
func Create(c *gin.Context) {
	b := newBuild()

	log.Debug("Create: save source to tmp")
	if err := b.saveSource(c, "./tmp"); err != nil {
		respondError(c, err)
		return
	}

	b.setup(c)

	log.Debug("Create: push build job")
	go jobqueue.Push(b.job)

	log.Debug("Create: success")
	c.JSON(http.StatusOK, gin.H{"id": b.job.Id})
}

func respondError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{"msg": err})
}

func newBuild() Build {
	return Build{job: jobqueue.NewJob()}
}

func save(r io.Reader, path string) error {
	w, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	if _, err := io.Copy(w, r); err != nil {
		return err
	}
	return nil
}

func formFileReader(c *gin.Context) (io.Reader, error) {
	if r, _, err := c.Request.FormFile("file"); err != nil {
		return nil, err
	} else {
		return r, nil
	}
}
