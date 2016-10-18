// +build !release

package env

import (
	log "github.com/Sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.Info("*** Debug mode is enabled ***")
}
