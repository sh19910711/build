package util

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
)

func Untar(r io.Reader, dstPrefix string) error {
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
		dstPath := filepath.Join(dstPrefix, header.Name)
		f, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE, os.FileMode(header.Mode))
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
