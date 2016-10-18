package archive

import (
	"archive/tar"
	"bytes"
	"io/ioutil"
	"path"
)

type Tar struct {
	Content []byte
	Name    string
	Mode    int64
	Size    int64
}

// TODO: options for tar header
func TarFromFile(filepath string) (*Tar, error) {
	if content, err := ioutil.ReadFile(filepath); err != nil {
		return nil, err
	} else {
		return &Tar{content, path.Base(filepath), 0644, int64(len(content))}, nil
	}
}

func TarFromBuffer(b *bytes.Buffer, filename string) *Tar {
	return &Tar{b.Bytes(), filename, 0644, int64(b.Len())}
}

func (t *Tar) Reader() (*bytes.Reader, error) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)

	if err := tw.WriteHeader(&tar.Header{Name: t.Name, Mode: t.Mode, Size: t.Size}); err != nil {
		return nil, err
	}
	if _, err := tw.Write(t.Content); err != nil {
		return nil, err
	}
	if err := tw.Close(); err != nil {
		return nil, err
	}

	return bytes.NewReader(buf.Bytes()), nil
}
