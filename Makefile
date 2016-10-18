.PHONY: build release init sync

build: sync
	go build -i -o bin/app

release: sync
	CGO_ENABLED=0 go build -a -installsuffix cgo -tags release -o bin/app

sync:
	govendor sync
	mkdir -p bin

init:
	go get -u github.com/kardianos/govendor
	govendor init
	govendor sync
	govendor install
