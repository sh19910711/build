.PHONY: build
build:
	govendor sync
	go build -i

run: build
	go run main.go

fmt:
	go fmt `go list ./... | grep -v \/vendor\/`

test:
	go test -v `go list ./... | grep -v \/vendor\/`

init:
	go get -u github.com/kardianos/govendor
	govendor init
