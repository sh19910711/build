package controller

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codestand/build/jobqueue"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
)

// POST /builds
// - params[file] := tar-ball (required)
// - params[callback] := URL fired after completed build (required)
// - returns {"id": "<job-id>"}
func Create(c *gin.Context) {
	log.Debug("1. create new job")
	job := jobqueue.NewJob()

	log.Debug("2. save source to tmp")
	r := formfile(c)
	if err := os.MkdirAll("./tmp", 0700); err != nil {
		log.Fatal(err)
	}
	src := "./tmp/" + job.Id
	if err := save(r, src); err != nil {
		log.Fatal(err)
	}

	log.Debug("3. push job to jobqueue")
	job.Src = src
	job.Callback = c.PostForm("callback")
	job.Image = "build"
	job.Commands = []string{"make"}
	job.Artifacts = []string{"/app/app"}
	jobqueue.Push(job)

	log.Debug("4. return job id")
	c.JSON(http.StatusOK, gin.H{"id": job.Id})
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

func formfile(c *gin.Context) io.Reader {
	r, _, err := c.Request.FormFile("file")
	if err != nil {
		log.Fatal("formfile: ", err)
	}
	return r
}
