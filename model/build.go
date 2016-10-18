package model

import (
	"github.com/jinzhu/gorm"
)

type Build struct {
	Id         int64
	SourceFile []byte
	Log        string
}

func All() []Build {
	builds := []Build{}
	db.Find(&builds)
	return builds
}

func Find(b *Build) *gorm.DB {
	return db.Find(b)
}
