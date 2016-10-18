package archive

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"io"
)

func ZipToTar(zipfile io.ReaderAt, len int64) (nilReader io.Reader, err error) {
	zr, err := zip.NewReader(zipfile, len)
	if err != nil {
		return nilReader, err
	}

	buf := &bytes.Buffer{}
	tw := tar.NewWriter(buf)

	err = func() error {
		for _, f := range zr.File {
			r, err := f.Open()
			if err != nil {
				return err
			}
			if err := tw.WriteHeader(&tar.Header{Name: f.Name, Mode: int64(f.Flags), Size: int64(f.UncompressedSize64)}); err != nil {
				return err
			}
			io.Copy(tw, r)
			r.Close()
		}
		if err := tw.Close(); err != nil {
			return err
		}
		return nil
	}()

	if err != nil {
		return nilReader, err
	} else {
		return bytes.NewReader(buf.Bytes()), nil
	}
}
