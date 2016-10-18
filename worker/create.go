package worker

import (
	"github.com/docker/docker/api/types/container"
	"golang.org/x/net/context"
)

func (w *Worker) Create(ctx context.Context) (err error) {
	config := container.Config{
		Image: w.Image,
		Cmd:   w.Cmd,
	}

	// Use command defined in Dockerfile if exists
	if info, _, err := w.c.ImageInspectWithRaw(ctx, w.Image); err != nil {
		return err
	} else if len(info.Config.Cmd) > 0 {
		config.Cmd = info.Config.Cmd
	}

	if c, err := w.c.ContainerCreate(ctx, &config, nil, nil, ""); err != nil {
		return err
	} else {
		w.Id = c.ID
	}

	return nil
}
