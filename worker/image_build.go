package worker

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/codestand/build/util"
	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
	"io"
	"io/ioutil"
	"strings"
)

type ImageBuildError struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type ImageBuildResponse struct {
	Stream      string           `json:"stream,omitempty"`
	ErrorDetail *ImageBuildError `json:"errorDetail,omitempty"`
}

func (w *Worker) ImageBuild(ctx context.Context, imageTag string, dockerfile io.Reader) error {
	// buildOptions can limit compute resources for builds
	options := types.ImageBuildOptions{}

	// archvie dockerfile
	b, err := ioutil.ReadAll(dockerfile)
	if err != nil {
		return err
	}
	r, err := util.ArchiveBuffer(bytes.NewBuffer(b), "Dockerfile")
	if err != nil {
		return err
	}

	// build image
	if res, err := w.c.ImageBuild(ctx, r, options); err != nil {
		return err
	} else {
		// read build log
		defer res.Body.Close()
		dec := json.NewDecoder(res.Body)

		for {
			var r ImageBuildResponse
			if err := dec.Decode(&r); err != nil {
				if err == io.EOF {
					break
				}
				return err
			}

			// TODO: improve log handling
			if r.ErrorDetail != nil {
				return errors.New(r.ErrorDetail.Message)
			} else {
				// save image
				var imageId string
				if strings.HasPrefix(r.Stream, "Successfully built") {
					fmt.Sscanf(r.Stream, "Successfully built %s", &imageId)
					if err := w.c.ImageTag(ctx, imageId, imageTag); err != nil {
						return err
					}
					return nil
				}
			}
		}

		return errors.New("build failed")
	}
}
