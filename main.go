package main

import (
	"github.com/codestand/build/controller"
	_ "github.com/codestand/build/env"
	"github.com/codestand/build/jobqueue"
)

func main() {
	go jobqueue.Wait()
	defer jobqueue.Close()
	controller.Mount()
}
