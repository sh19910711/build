package testhelper

import (
	"bytes"
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
)

func Init() {
	log.SetLevel(log.DebugLevel)
}

func UploadRequest(url string, key string, path string, params map[string]string) (req *http.Request, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	part, err := w.CreateFormFile(key, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	io.Copy(part, file)

	for k, v := range params {
		w.WriteField(k, v)
	}
	if err := w.Close(); err != nil {
		return nil, err
	}

	req, err = http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	return req, nil
}

func Serve(path string, handler func(c *gin.Context)) *httptest.Server {
	r := gin.Default()
	r.POST(path, handler)
	return httptest.NewServer(r)
}
