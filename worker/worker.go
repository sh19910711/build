package worker

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
	"io"
	"os"
	"strings"
)

const DOCKER_ENDPOINT = "unix:///var/run/docker.sock"
const DOCKER_API_VERSION = "v1.18"

func New() (w Worker, err error) {
	headers := map[string]string{
		"User-Agent": "engine-api",
	}
	endpoint := os.Getenv("DOCKER_HOST")
	if endpoint == "" {
		endpoint = DOCKER_ENDPOINT
	}
	w.c, err = client.NewClient(endpoint, DOCKER_API_VERSION, nil, headers)
	return w, err
}

type Worker struct {
	id string
	c  *client.Client // docker engine client
}

func (w *Worker) Create(ctx context.Context, imageName string, cmd string) (err error) {
	config := container.Config{
		Image: imageName,
		Cmd:   strings.Split(cmd, " "),
	}

	c, err := w.c.ContainerCreate(ctx, &config, nil, nil, "")
	if err != nil {
		return err
	}

	w.id = c.ID
	return nil
}

func (w *Worker) CopyToWorker(ctx context.Context, tar io.Reader, dst string) error {
	opts := types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: true,
	}
	return w.c.CopyToContainer(ctx, w.id, dst, tar, opts)
}

// Get file from Worker (as a tar-ball archive).
func (w *Worker) CopyFromWorker(ctx context.Context, file string) (io.Reader, error) {
	r, _, err := w.c.CopyFromContainer(ctx, w.id, file)
	return r, err
}

func (w *Worker) Start(ctx context.Context) error {
	return w.c.ContainerStart(ctx, w.id, types.ContainerStartOptions{})
}

func (w *Worker) Wait(ctx context.Context) (int, error) {
	return w.c.ContainerWait(ctx, w.id)
}

func (w *Worker) Destroy(ctx context.Context) error {
	return w.c.ContainerRemove(ctx, w.id, types.ContainerRemoveOptions{})
}
