package model

import (
	"bytes"
	"github.com/codestand/build/archive"
	"io"
)

func (b *Build) AppTar() (io.Reader, error) {
	br := bytes.NewReader(b.SourceFile)
	return archive.ZipToTar(br, int64(br.Len()))
}
