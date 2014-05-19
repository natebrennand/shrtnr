
EXECUTABLE=shrtnr

default: build

test:
	go test ./...

build:
	go build -o $(EXECUTABLE) main.go routes.go

clean:
	go clean

.phony: test build clean
