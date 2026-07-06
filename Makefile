BINARY := statusline
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -s -w -X main.version=$(VERSION)

.PHONY: build test test-integration test-e2e bench fmt lint install clean

build:
	go build -trimpath -ldflags="$(LDFLAGS)" -o $(BINARY) ./cmd/statusline

test:
	go test ./...

test-integration:
	go test -tags integration ./...

test-e2e:
	go test -tags e2e ./...

bench:
	go test -bench=. -benchmem -run=^$$ ./...

fmt:
	gofmt -s -l -w .

lint:
	go vet ./...
	golangci-lint run ./...

install:
	go install -trimpath -ldflags="$(LDFLAGS)" ./cmd/statusline

clean:
	rm -f $(BINARY)
