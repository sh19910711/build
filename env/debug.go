// +build !release

package env

import (
	log "github.com/Sirupsen/logrus"
)

const DEBUG = true

func Init() {
	log.SetLevel(log.DebugLevel)
	log.Info("*** Debug mode is enabled ***")
}
