package worker

import (
	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
	"io"
)

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
