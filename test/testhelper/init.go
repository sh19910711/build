package testhelper

import (
	"flag"
	"os"
)

func Init() {
	cwd := flag.String("cwd", "", "set cwd")
	flag.Parse()
	if *cwd != "" {
		if err := os.Chdir(*cwd); err != nil {
			println(err)
		}
	}
}
