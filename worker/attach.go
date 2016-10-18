package worker

import (
	"encoding/binary"
	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
	"io"
)

type DockerStreamHeader struct {
	Type     byte
	Reserved [3]byte
	Size     uint32
}

func (w *Worker) Attach(ctx context.Context) (r io.Reader, err error) {
	opts := types.ContainerAttachOptions{Stream: true, Stdin: false, Stdout: true, Stderr: true}
	resp, err := w.c.ContainerAttach(ctx, w.Id, opts)
	if err != nil {
		return r, err
	}
	resp.CloseWrite()

	in, out := io.Pipe()
	go func() {
		for {
			// stream format: header => payload
			h := DockerStreamHeader{}
			if err := binary.Read(resp.Reader, binary.BigEndian, &h); err == io.EOF {
				break
			} else if err != nil {
				panic(err)
			}
			io.CopyN(out, resp.Reader, int64(h.Size))
		}

		out.Close()
	}()

	return in, nil
}
