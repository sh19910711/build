package worker_test

import (
	"golang.org/x/net/context"
	"time"
)

func contextWithTimeout() (context.Context, func()) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}
