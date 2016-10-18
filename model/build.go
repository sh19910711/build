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
}

func All() []Build {
	builds := []Build{}
	db.Find(&builds)
	return builds
}

func Create(id int64) Build {
	b := Build{}
	db.Create(&b)
	return b
}

func Find(b *Build) *gorm.DB {
	return db.Find(b)
}

func (b *Build) WriteLog(msg string) *gorm.DB {
	return db.Model(b).UpdateColumn("log", b.Log+"\n"+msg)
}
