package worker

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
	"io"
	"os"
)

const DOCKER_ENDPOINT = "unix:///var/run/docker.sock"
const DOCKER_API_VERSION = "v1.18"

var dockerClient *client.Client

func init() {
	headers := map[string]string{
		"User-Agent": "engine-api",
	}
	endpoint := os.Getenv("DOCKER_HOST")
	if endpoint == "" {
		endpoint = DOCKER_ENDPOINT
	}
	if cli, err := client.NewClient(endpoint, DOCKER_API_VERSION, nil, headers); err != nil {
		panic(err)
	} else {
		dockerClient = cli
	}
}

type Worker struct {
	Id    string
	c     *client.Client // docker engine client
	Image string
	Cmd   []string
}

func New() (w Worker) {
	return Worker{
		c:     dockerClient,
		Image: "build",
		Cmd:   []string{"bash", "/build.bash"},
	}
}

func (w *Worker) Create(ctx context.Context) (err error) {
	config := container.Config{
		Image: w.Image,
		Cmd:   w.Cmd,
	}

	// use image command if exist for now
	if info, _, err := w.c.ImageInspectWithRaw(ctx, w.Image); err != nil {
		return err
	} else {
		if len(info.Config.Cmd) > 0 {
			config.Cmd = info.Config.Cmd
		}
	}

	c, err := w.c.ContainerCreate(ctx, &config, nil, nil, "")
	if err != nil {
		return err
	}

	w.Id = c.ID
	return nil
}

func (w *Worker) CopyToWorker(ctx context.Context, tar io.Reader, dst string) error {
	opts := types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: true,
	}
	return w.c.CopyToContainer(ctx, w.Id, dst, tar, opts)
}

// Get file from Worker (as a tar-ball archive).
func (w *Worker) CopyFromWorker(ctx context.Context, file string) (io.Reader, error) {
	r, _, err := w.c.CopyFromContainer(ctx, w.Id, file)
	return r, err
}

func (w *Worker) Start(ctx context.Context) error {
	return w.c.ContainerStart(ctx, w.Id, types.ContainerStartOptions{})
}

func (w *Worker) Wait(ctx context.Context) (int, error) {
	return w.c.ContainerWait(ctx, w.Id)
}

func (w *Worker) Destroy(ctx context.Context) error {
	return w.c.ContainerRemove(ctx, w.Id, types.ContainerRemoveOptions{})
}

func (w *Worker) IsFinished(ctx context.Context) (bool, error) {
	c, err := w.c.ContainerInspect(ctx, w.Id)
	if err != nil {
		return false, err
	}
	return c.State.Status == "exited", nil
}
