package archive

import (
	"archive/tar"
	"bytes"
	"io"
	"io/ioutil"
	"path"
)

type Tar struct {
	Content []byte
	Name    string
	Mode    int64
	Size    int64
}

// TODO: file path?
func GetFileReaderFromTar(tarball io.Reader, filename string) (nilReader io.Reader, err error) {
	tr := tar.NewReader(tarball)

	for { // each file
		header, err := tr.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return nilReader, err
		}

		if header.Name != filename {
			continue
		}

		buf := &bytes.Buffer{}
		if _, err := io.CopyN(buf, tr, header.Size); err != nil {
			return nilReader, err
		}
		return bytes.NewReader(buf.Bytes()), nil
	}

	return nilReader, nil
}

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
	defer tw.Close()

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
