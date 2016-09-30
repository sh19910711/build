.PHONY: build release run test test_all init archive

build:
	govendor sync
	go build -i -o app

release:
	govendor sync
	CGO_ENABLED=0 go build -a -installsuffix cgo -tags release -o app

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
	sudo chmod -R 0755 app
	sudo chown -R root. app
	tar zcvf app.tar.gz app
