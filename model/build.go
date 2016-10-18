package model

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
