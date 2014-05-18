
EXECUTABLE=shrtnr

default: build

run:
	redis-server ./config/redis.conf
	$(EXECUTABLE)

test:
	go test ./...

build:
	go build -o $(EXECUTABLE) main.go

clean:
	go clean

.phony: test build clean
