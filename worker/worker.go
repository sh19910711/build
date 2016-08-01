package worker

import(
	log "github.com/Sirupsen/logrus"
	"bytes"
	"archive/tar"
	"io/ioutil"
	"path"
	"io"
	"os"
	"golang.org/x/net/context"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
)

const DOCKER_ENDPOINT = "unix:///var/run/docker.sock"
const DOCKER_API_VERSION = "v1.18"

func New() (w worker, err error) {
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

func (w *worker) CopyFromWorker(src, dstPrefix string) error {
	// mkdir dstPrefix
	if err := os.MkdirAll(dstPrefix, 0755); err != nil {
		return err
	}
	// get file from worker (as a tar-ball archive)
	r, _, err := w.c.CopyFromContainer(w.ctx, w.id, src)
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
		log.Info("artifacts: ", header.Name)
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

func (w *worker) Start() error {
	return w.c.ContainerStart(w.ctx, w.id, types.ContainerStartOptions{})
}

func (w *worker) Wait() (int, error) {
	return w.c.ContainerWait(w.ctx, w.id)
}


