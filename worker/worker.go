package worker

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

type Worker struct {
	Id    string
	c     *client.Client
	Image string
	Cmd   []string
}

func New() *Worker {
	return &Worker{c: dockerClient}
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
