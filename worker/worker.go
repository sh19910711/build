package worker

import (
	"archive/tar"
	log "github.com/Sirupsen/logrus"
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

func (w *Worker) Copy(ctx context.Context, tar io.Reader, dst string) error {
	opts := types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: true,
	}
	return w.c.CopyToContainer(ctx, w.id, dst, tar, opts)
}

func (w *Worker) CopyFile(ctx context.Context, src string, dst string) error {
	r, err := archive(src)
	if err != nil {
		return err
	}
	return w.Copy(ctx, r, dst)
}

func untar(r io.Reader, dstPrefix string) error {
	// extract artifacts from archive
	tr := tar.NewReader(r)

	// iterate through the files
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// write
		log.Debug("artifacts: ", header.Name)
		f, err := os.OpenFile(dstPrefix+"/"+header.Name, os.O_WRONLY|os.O_CREATE, os.FileMode(header.Mode))
		defer f.Close()
		if err != nil {
			return err
		}
		if _, err := io.Copy(f, tr); err != nil {
			return err
		}
	}

	return nil
}

func (w *Worker) CopyFromWorker(ctx context.Context, file string) (io.Reader, error) {
	// get file from Worker (as a tar-ball archive)
	r, _, err := w.c.CopyFromContainer(ctx, w.id, file)
	if err != nil {
		return nil, err
	}

	return r, nil
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
