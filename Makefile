M = $(shell printf "\033[34;1m▶\033[0m")

all: build

dep: ; $(info $(M) Ensuring dependencies)
	go get github.com/ahmetb/govvv
	GO111MODULE=on go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.15.0

clean: ; $(info $(M) Cleaning)
	rm -f homlet

test: clean dep ; $(info $(M) Launching tests)
	go test $(TEST_FLAGS) ./...

fmt: ; $(info $(M) Formatting code)
	gofmt -s -w .

lint: ; $(info $(M) Linting)
	golangci-lint run

build: clean ; $(info $(M) Building)
	go build -ldflags "$(shell govvv -flags -pkg $(shell go list ./pkg/version)) -s -w" ./cmd/homlet

rpi: clean ; $(info $(M) Building raspberry pi binary)
	GOOS=linux GOARCH=arm GOARM=5 go build -ldflags "$(shell govvv -flags -pkg $(shell go list ./pkg/version)) -s -w" ./cmd/homlet

release: lint build
