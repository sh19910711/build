package testhelper

import (
	"archive/tar"
	"io"
)

func ShouldIncludeFileInTar(r io.Reader, path string) (bool, error) {
	tr := tar.NewReader(r)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return false, err
		}
		if header.Name == path {
			return true, nil
		}
	}

	return false, nil
}
