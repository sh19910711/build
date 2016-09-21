package util

import (
	"archive/tar"
	"bytes"
	"io/ioutil"
	"path"
)

func ArchiveFile(src string) (*bytes.Reader, error) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)

	content, err := ioutil.ReadFile(src)
	if err != nil {
		return nil, err
	}

	h := &tar.Header{
		Name: path.Base(src),
		Mode: 0755,
		Size: int64(len(content)),
	}
	if err := tw.WriteHeader(h); err != nil {
		return nil, err
	}
	if _, err := tw.Write(content); err != nil {
		return nil, err
	}
	if err := tw.Close(); err != nil {
		return nil, err
	}

	return bytes.NewReader(buf.Bytes()), nil
}

func ArchiveBuffer(in *bytes.Buffer, filename string) (*bytes.Reader, error) {
	out := new(bytes.Buffer)
	tw := tar.NewWriter(out)

	h := &tar.Header{
		Name: filename,
		Mode: 0644,
		Size: int64(in.Len()),
	}
	if err := tw.WriteHeader(h); err != nil {
		return nil, err
	}
	if _, err := tw.Write(in.Bytes()); err != nil {
		return nil, err
	}
	if err := tw.Close(); err != nil {
		return nil, err
	}

	return bytes.NewReader(out.Bytes()), nil
}
