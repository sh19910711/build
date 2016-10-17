.PHONY: build release init

build:
	govendor sync
	go build -i -o bin/app

release:
	govendor sync
	CGO_ENABLED=0 go build -a -installsuffix cgo -tags release -o bin/app

init:
	go get -u github.com/kardianos/govendor
	govendor init
	govendor sync
