all: run

run: install
	oakblue

fmt:
	go fmt ./...

test: fmt
	go test -count=1 ./...

install: fmt
	go install ./...