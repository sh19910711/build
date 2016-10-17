package main

import (
	"github.com/codestand/build/controller"
	_ "github.com/codestand/build/env"
)

func main() {
	controller.Mount()
}
