package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codestand/build/worker"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"msg": "hello"})
	})

	r.POST("/builds", func(c *gin.Context) {
		// new worker
		w, err := worker.New()
		if err != nil {
			log.Fatal(err)
		}
		log.Info("worker has been initialized")

		// create a worker
		if err := w.Create(); err != nil {
			log.Fatal(err)
		}
		log.Info("worker has been created")

		// send script
		if err := w.CopyFile("./script/build.bash", "/"); err != nil {
			log.Fatal(err)
		}
		log.Info("script has been sent to worker")

		// send app
		file, _, err := c.Request.FormFile("f")
		if err != nil {
			log.Fatal(err)
		}
		if err := w.Copy(file, "/app"); err != nil {
			log.Fatal(err)
		}
		log.Info("app has been sent to worker")

		// start a worker
		if err := w.Start(); err != nil {
			log.Fatal(err)
		}
		log.Info("worker has been started")

		// wait a worker
		exitCode, err := w.Wait()
		if err != nil {
			log.Fatal(err)
		}
		log.Info("worker has been exited with ", exitCode)

		// test to get an artifact
		if err := w.CopyFromWorker("/app/app", "./tmp"); err != nil {
			log.Fatal(err)
		}
		log.Info("success build")
	})

	r.Run()
}
