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

func archiveDockerfile(in io.Reader) (nilReader io.Reader, err error) {
	b, err := ioutil.ReadAll(in)
	if err != nil {
		return nilReader, err
	}

	if r, err := util.ArchiveBuffer(bytes.NewBuffer(b), "Dockerfile"); err != nil {
		return nilReader, err
	} else {
		return r, nil
	}
}

func getImageIdFromResponseBody(resBody io.Reader) (nilImageId string, err error) {
	dec := json.NewDecoder(resBody)

	for { // each command and its output
		var r ImageBuildResponse
		if err := dec.Decode(&r); err != nil {
			if err == io.EOF {
				break
			}
			return nilImageId, err
		}

		if r.ErrorDetail != nil {
			return nilImageId, errors.New(r.ErrorDetail.Message)
		} else {
			if strings.HasPrefix(r.Stream, "Successfully built") {
				var imageId string
				fmt.Sscanf(r.Stream, "Successfully built %s", &imageId)
				return imageId, nil
			}
		}
	}

	return nilImageId, errors.New("build failed")
}

func (w *Worker) ImageBuild(ctx context.Context, dockerfile io.Reader) error {
	// the options can limit compute resources for builds
	options := types.ImageBuildOptions{}

	r, err := archiveDockerfile(dockerfile)
	if err != nil {
		return err
	}

	if res, err := w.c.ImageBuild(ctx, r, options); err != nil {
		return err
	} else {
		defer res.Body.Close()
		if imageId, err := getImageIdFromResponseBody(res.Body); err != nil {
			return err
		} else {
			if err := w.c.ImageTag(ctx, imageId, w.Image); err != nil {
				return err
			}
			return nil
		}
	}
}
