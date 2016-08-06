package worker

import (
	"archive/tar"
	"bytes"
	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"golang.org/x/net/context"
	"io"
	"io/ioutil"
	"os"
	"path"
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
	ID string
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

	w.ID = c.ID
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

func (w *Worker) Copy(ctx context.Context, tar io.Reader, dst string) error {
	opts := types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: true,
	}
	return w.c.CopyToContainer(ctx, w.ID, dst, tar, opts)
}

func (w *Worker) CopyFile(ctx context.Context, src string, dst string) error {
	r, err := archive(src)
	if err != nil {
		return err
	}
	return w.Copy(ctx, r, dst)
}

func (w *Worker) CopyFromWorker(ctx context.Context, src, dstPrefix string) error {
	// mkdir dstPrefix
	if err := os.MkdirAll(dstPrefix, 0755); err != nil {
		return err
	}
	// get file from Worker (as a tar-ball archive)
	r, _, err := w.c.CopyFromContainer(ctx, w.ID, src)
	if err != nil {
		return err
	}

	// extract artifacts from archive
	tr := tar.NewReader(r)

	// iterate through the files
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		// write
		log.Debug("artifacts: ", header.Name)
		f, err := os.OpenFile(dstPrefix+"/"+header.Name, os.O_WRONLY|os.O_CREATE, os.FileMode(header.Mode))
		defer f.Close()
		if err != nil {
			log.Fatal(err)
		}
		if _, err := io.Copy(f, tr); err != nil {
			log.Fatal(err)
		}
	}

	return nil
}

func (w *Worker) Start(ctx context.Context) error {
	return w.c.ContainerStart(ctx, w.ID, types.ContainerStartOptions{})
}

func (w *Worker) Wait(ctx context.Context) (int, error) {
	return w.c.ContainerWait(ctx, w.ID)
}

func (w *Worker) Destroy(ctx context.Context) error {
	return w.c.ContainerRemove(ctx, w.ID, types.ContainerRemoveOptions{})
}
