package main

import (
	"github.com/codestand/build/controller"
	_ "github.com/codestand/build/env"
	"github.com/codestand/build/job"
	"github.com/codestand/build/model"
)

func main() {
	model.Open()
	defer model.Close()
	go job.Wait()
	controller.Mount()
}
