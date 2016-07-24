package main

import log "github.com/Sirupsen/logrus"
import "github.com/gin-gonic/gin"
import "net/http"
import "io"
import "io/ioutil"
import "os"
import "archive/tar"
import "bytes"
import "github.com/docker/engine-api/client"
import "github.com/docker/engine-api/types"
import "github.com/docker/engine-api/types/container"
import "golang.org/x/net/context"

func main() {
	log.Info("starting build server")

	defaultHeaders := map[string]string{"User-Agent": "engine-api"}
	client, err := client.NewClient("unix:///var/run/docker.sock", "v1.18", nil, defaultHeaders)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"msg": "hello"})
	})

	r.GET("/docker/containers", func(c *gin.Context) {
		options := types.ContainerListOptions{All: true}
		containers, err := client.ContainerList(context.Background(), options)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{"containers": containers})
	})

  // exec test
  r.GET("/docker/exec", func(c *gin.Context) {
    // create container
    config := container.Config{
      Image: "ubuntu:16.04",
      Cmd: []string{"sleep", "3"},
    }
    worker, err := client.ContainerCreate(context.Background(), &config, nil, nil, "")
    if err != nil {
      log.Fatal(err)
    }
    log.Info(worker.ID)

    // exec
    ctx := context.Background()
    ec := types.ExecConfig{
      Cmd: []string{"echo", "hello"},
    }
    res, err := client.ContainerExecCreate(ctx, worker.ID, ec)
    if err != nil {
      log.Fatal(err)
    }
    execID := res.ID
    if execID == "" {
      log.Fatal("exec id empty")
    }
  })

	r.POST("/tar", func(c *gin.Context) {
		// get file
		file, header, err := c.Request.FormFile("f")
		if err != nil {
			log.Fatal(err)
		}
		log.Info("filename: ", header.Filename)

		// create tar ball reader
		buf := bytes.NewBuffer(nil)
		io.Copy(buf, file)
		tr := tar.NewReader(buf)

		// read files
		for {
			hdr, err := tr.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			log.Infof("Contents of %s:\n", hdr.Name)
			if _, err := io.Copy(os.Stdout, tr); err != nil {
				log.Fatal(err)
			}
		}
	})

	r.POST("/upload", func(c *gin.Context) {
		// get file
		file, header, err := c.Request.FormFile("f")
		if err != nil {
			log.Fatalln(err)
		}
		filename := header.Filename
		log.Info("filename: ", filename)

		// create tmpdir
		tmpdir, err := ioutil.TempDir("", "build")
		if err != nil {
			log.Fatalln(err)
		}
		path := tmpdir + "/" + filename
		log.Info("saved into " + path)
		out, err := os.Create(path)
		defer out.Close()

		// save into tmpdir
		_, err = io.Copy(out, file)
		if err != nil {
			log.Fatalln(err)
		}

		c.JSON(http.StatusOK, gin.H{"msg": "uploaded"})
	})

	r.Run()
}
