build:
	govendor sync
	go build -i

run:
	govendor sync
	go run main.go

fmt:
	go fmt ./...

init:
	go get -u github.com/kardianos/govendor
	govendor init
