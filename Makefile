.PHONY: build release init sync test

build: sync
	go build -i -o bin/app

release: sync
	CGO_ENABLED=0 go build -a -installsuffix cgo -tags release -o bin/app

sync:
	govendor sync
	mkdir -p bin

test: sync
	go test -v `go list ./... | grep -v "/vendor/"`

init:
	go get -u github.com/kardianos/govendor
	govendor init
	govendor sync
	govendor install
