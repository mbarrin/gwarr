test:
	go test ./...

build:
	go build -o dist/ ./cmd/gwarr/gwarr.go

build-docker: build
	docker build -t gwarr:local .

clean:
	rm -rf build/

clean-docker:
	docker rmi gwarr:local

clean-all: clean clean-docker

lint:
	golangci-lint run
