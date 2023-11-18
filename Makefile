test:
	go test ./...

build:
	go build -o build/ ./cmd/gwarr/gwarr.go

build-docker: build
	docker build -t gwarr:latest .

clean:
	rm -rf build/

clean-docker:
	docker rmi gwarr:latest

clean-all: clean clean-docker

lint:
	golangci-lint run
