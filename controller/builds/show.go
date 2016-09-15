package builds

import (
	"github.com/codestand/build/job"
	"github.com/codestand/build/worker"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GET /builds/<build-id>
func Show(c *gin.Context) {
	id := c.Param("id")
	if j, err := job.Find(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err})
	} else {
		finished, err := worker.IsFinished(j.WorkerId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": err})
		} else {
			c.JSON(http.StatusOK, gin.H{"id": id, "finished": finished, "exitCode": j.ExitCode})
		}
	}
}
