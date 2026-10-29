PACKAGE=github.com/genshen/ssh-web-console
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS=-ldflags "-X main.buildVersion=${VERSION}"

.PHONY: clean all

all: ssh-web-console-linux-amd64 ssh-web-console-linux-arm64 ssh-web-console-darwin-amd64 ssh-web-console-windows-amd64.exe ssh-web-console-freebsd-amd64 ssh-web-console-freebsd-arm64

ssh-web-console-linux-amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ssh-web-console-linux-amd64 ${PACKAGE}

ssh-web-console-linux-arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o ssh-web-console-linux-arm64 ${PACKAGE}

ssh-web-console-darwin-amd64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o ssh-web-console-darwin-amd64 ${PACKAGE}

ssh-web-console-windows-amd64.exe:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o ssh-web-console-windows-amd64.exe ${PACKAGE}

ssh-web-console-freebsd-amd64:
	CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 go build ${LDFLAGS} -o ssh-web-console-freebsd-amd64 ${PACKAGE}

ssh-web-console-freebsd-arm64:
	CGO_ENABLED=0 GOOS=freebsd GOARCH=arm64 go build ${LDFLAGS} -o ssh-web-console-freebsd-arm64 ${PACKAGE}

ssh-web-console :
	go build ${LDFLAGS} -o ssh-web-console

clean:
	rm -f ssh-web-console-linux-amd64 ssh-web-console-linux-arm64 ssh-web-console-darwin-amd64 ssh-web-console-windows-amd64.exe ssh-web-console-freebsd-amd64 ssh-web-console-freebsd-arm64
