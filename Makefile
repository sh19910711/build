.PHONY: build run test test_all init

build:
	govendor sync
	go build -i -o app

release:
  govendor sync
  go build -tags release -i -o app

run: build
	go run main.go

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

archive: release
	tar zcvf app.tar.gz app
