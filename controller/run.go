package controller

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codestand/build/job"
	"github.com/codestand/build/model"
	"github.com/gin-gonic/gin"
)

func Run(c *gin.Context) {
	b := &model.Build{}

	if id, err := atoi(c.Param("id")); err != nil {
		log.Warn(err)
		internalError(c)
		return
	} else {
		b.Id = id
	}

	if model.Find(b).RecordNotFound() {
		notFound(c)
	} else {
		job.Push(c.Param("id"))
		jsonResponse(c, gin.H{"msg": "ok"})
	}
}
