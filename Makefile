DEFAULT_VERSION=0.1.0-local
VERSION := $(or $(VERSION),$(DEFAULT_VERSION))

all: binaries
run:
	go run -ldflags "-w -X main.version=$(VERSION)" main.go
clean:
	rm -rf bin
test:
	go test -v ./...
binaries:
	gox \
		-osarch="linux/amd64 linux/arm64 darwin/amd64 windows/amd64" \
		-ldflags "-w -X main.version=$(VERSION)" \
		-output="bin/{{.OS}}/{{.Arch}}/app" \
		./...