package model_test

import (
	"os"
)

func init() {
	if os.Getenv("DATABASE_URL") == "" {
		panic("DATABASE_URL must be set")
	}
}
