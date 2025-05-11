.PHONY: build clean test release snapshot

VERSION := 0.1.2
BUILD_DATE := $(shell date +%Y-%m-%d)
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildDate=$(BUILD_DATE)"

build:
	go build $(LDFLAGS) -o bin/term-rex

clean:
	rm -rf bin/ dist/

test:
	go test ./...

vet:
	go vet ./...

# Local build for all platforms
build-all:
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/term-rex-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/term-rex-darwin-arm64
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/term-rex-linux-amd64
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o bin/term-rex-linux-arm64
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/term-rex-windows-amd64.exe

# Create a new release using goreleaser
release:
	goreleaser release --clean

# Test the release process without publishing
snapshot:
	goreleaser release --snapshot --clean --skip=publish

# Create a new git tag for release
tag:
	git tag -a v$(VERSION) -m "Release v$(VERSION)"
	git push origin v$(VERSION)
