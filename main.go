package main

import (
	"archive/tar"
	"bytes"
	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	_ "time"
)

const DOCKER_ENDPOINT = "unix:///var/run/docker.sock"
const DOCKER_API_VERSION = "v1.18"

func newWorker() (w worker, err error) {
	headers := map[string]string{
		"User-Agent": "engine-api",
	}
	w.c, err = client.NewClient(DOCKER_ENDPOINT, DOCKER_API_VERSION, nil, headers)
	return w, err
}

type worker struct {
	c   *client.Client
	id  string
	ctx context.Context
}

func (w *worker) Create() (err error) {
	config := container.Config{
		Image: "build",
		Cmd:   []string{"bash", "/build.bash"},
	}

	w.ctx = context.Background()
	c, err := w.c.ContainerCreate(w.ctx, &config, nil, nil, "")
	if err != nil {
		return err
	}

	w.id = c.ID
	return nil
}

func archive(src string) (*bytes.Reader, error) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)

	content, err := ioutil.ReadFile(src)
	if err != nil {
		return nil, err
	}

	h := &tar.Header{
		Name: path.Base(src),
		Mode: 0755,
		Size: int64(len(content)),
	}
	if err := tw.WriteHeader(h); err != nil {
		return nil, err
	}
	if _, err := tw.Write(content); err != nil {
		return nil, err
	}
	if err := tw.Close(); err != nil {
		return nil, err
	}

	return bytes.NewReader(buf.Bytes()), nil
}

func (w *worker) Copy(r io.Reader, dst string) error {
	opts := types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: true,
	}
	return w.c.CopyToContainer(w.ctx, w.id, dst, r, opts)
}

func (w *worker) CopyFile(src string, dst string) error {
	r, err := archive(src)
	if err != nil {
		return err
	}
	return w.Copy(r, dst)
}

func (w *worker) Start() error {
	return w.c.ContainerStart(w.ctx, w.id, types.ContainerStartOptions{})
}

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"msg": "hello"})
	})

	r.POST("/builds", func(c *gin.Context) {
		// new worker
		w, err := newWorker()
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
			log.Fatal("script: ", err)
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
	})

	r.Run()
}
