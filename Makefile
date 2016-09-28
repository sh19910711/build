.PHONY: build run test init

build:
	govendor sync
	go build -i -o app

run: build
	go run main.go

fmt:
	go fmt `go list ./... | grep -v \/vendor\/`

test:
	mkdir -p tmp
	go test -v `go list ./... | grep -v \/vendor\/`

test_all:
	mkdir -p tmp
	go test -tags integration -v `go list ./... | grep -v \/vendor\/` -cwd=$(PWD)

init:
	go get -u github.com/kardianos/govendor
	govendor init
	govendor sync
	govendor fetch +missing
