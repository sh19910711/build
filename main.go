package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codestand/build/worker"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"net/http"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"msg": "hello"})
	})

	r.POST("/builds", func(c *gin.Context) {
		ctx := context.Background()

		// new worker
		w, err := worker.New()
		if err != nil {
			log.Fatal(err)
		}
		log.Debug("worker has been initialized")

		// create a worker
		if err := w.Create(ctx, "build", "bash /build.bash"); err != nil {
			log.Fatal(err)
		}
		log.Debug("worker has been created")

		// send script
		if err := w.CopyFile(ctx, "./script/build.bash", "/"); err != nil {
			log.Fatal(err)
		}
		log.Debug("script has been sent to worker")

		// send app
		file, _, err := c.Request.FormFile("f")
		if err != nil {
			log.Fatal(err)
		}
		if err := w.Copy(ctx, file, "/app"); err != nil {
			log.Fatal(err)
		}
		log.Debug("app has been sent to worker")

		// start a worker
		if err := w.Start(ctx); err != nil {
			log.Fatal(err)
		}
		log.Debug("worker has been started")

		// wait a worker
		exitCode, err := w.Wait(ctx)
		if err != nil {
			log.Fatal(err)
		}
		log.Debug("worker has been exited with ", exitCode)

		// test to get an artifact
		if err := w.CopyFromWorker(ctx, "/app/app", "./tmp"); err != nil {
			log.Fatal(err)
		}
		log.Debug("success build")
	})

	r.Run()
}
