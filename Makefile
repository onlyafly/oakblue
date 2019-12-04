all: run

run: install
	oakblue

fmt:
	go fmt ./...

test: fmt
	go test ./...

install: fmt
	go install ./...