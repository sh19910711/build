package jobqueue

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codestand/build/util"
	"github.com/codestand/build/worker"
	"golang.org/x/net/context"
	"os"
)

func spawnJob(running job) {
	ctx := context.Background()

	// new worker
	w, err := worker.New()
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("jobqueue: worker initialized")

	// create a worker
	if err := w.Create(ctx, "build", "bash /build.bash"); err != nil {
		log.Fatal(err)
	}
	log.Debug("jobqueue: worker has been created")

	// send script
	if err := w.CopyFile(ctx, "./script/build.bash", "/"); err != nil {
		log.Fatal(err)
	}
	log.Debug("jobqueue: script has been sent to worker")

	// send app
	r, err := os.Open(running.Src)
	if err != nil {
		log.Fatal(err)
	}
	if err := w.Copy(ctx, r, "/app"); err != nil {
		log.Fatal(err)
	}
	log.Debug("jobqueue: app has been sent to worker")

	// start a worker
	if err := w.Start(ctx); err != nil {
		log.Fatal(err)
	}
	log.Debug("jobqueue: worker has been started")

	// wait a worker
	exitCode, err := w.Wait(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("jobqueue: worker has been exited with ", exitCode)

	// send artifact to callback url
	if running.Callback != "" {
		artifacts, err := w.CopyFromWorker(ctx, "/app/app")
		if err != nil {
			log.Fatal(err)
		}
		log.Debug("jobqueue: artifacts is found")

		if _, err := util.Upload(artifacts, running.Callback, "file", "artifacts.tar"); err != nil {
			log.Fatal(err)
		} else {
			log.Debug("jobqueue: fired callback")
		}
	}

	log.Debug("jobqueue: success :-)")
}
