package util

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
)

func CheckFileInTar(r io.Reader, filename string) (bool, error) {
	tr := tar.NewReader(r)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return false, err
		}
		if hdr.Name == filename {
			return true, nil
		}
	}
	return false, nil
}

func ReadFileFromTar(r io.Reader, filename string) (res io.Reader, err error) {
	pr, pw := io.Pipe()
	tr := tar.NewReader(r)

	// iterate through the files
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return res, err
		}

		if header.Name == filename {
			if _, err := io.Copy(pw, tr); err != nil {
				return res, err
			}
			return pr, nil
		}
	}

	return res, nil
}

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
