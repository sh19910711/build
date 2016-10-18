package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Build struct {
	Id         int64
	SourceFile []byte
	Log        string
	UpdatedAt  time.Time
	ImageName  string // the name of the docker image
}

func All() []Build {
	builds := []Build{}
	db.Find(&builds)
	return builds
}

func Create(id int64) Build {
	b := Build{ImageName: "codestand/baseos"}
	db.Create(&b)
	return b
}

func Find(b *Build) *gorm.DB {
	return db.Find(b)
}

func (b *Build) WriteLog(msg string) *gorm.DB {
	if b.Log == "" {
		b.Log = msg
	} else {
		b.Log += "\n" + msg
	}

	return db.Model(b).UpdateColumn("log", b.Log)
}

func (b *Build) ResetLog() *gorm.DB {
	b.Log = ""
	return db.Model(b).UpdateColumn("log", b.Log)
}
