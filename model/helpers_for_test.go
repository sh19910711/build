package model_test

import (
	"github.com/codestand/build/model"
)

func init() {
	model.Open()
	defer model.Close()
	if model.Find(&model.Build{Id: 10000}).RecordNotFound() {
		model.Create(10000)
	}
}
