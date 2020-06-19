PACKAGE=github.com/genshen/ssh-web-console

.PHONY: clean all

all: ssh-web-console-linux-amd64 ssh-web-console-linux-arm64 ssh-web-console-darwin-amd64 ssh-web-console-windows-amd64.exe

ssh-web-console-linux-amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ssh-web-console-linux-amd64 ${PACKAGE}

ssh-web-console-linux-arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ssh-web-console-linux-arm64 ${PACKAGE}

ssh-web-console-darwin-amd64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ssh-web-console-darwin-amd64 ${PACKAGE}

ssh-web-console-windows-amd64.exe:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ssh-web-console-windows-amd64.exe ${PACKAGE}

ssh-web-console :
	go build -o ssh-web-console

clean:
	rm -f ssh-web-console-linux-amd64 ssh-web-console-linux-arm64 ssh-web-console-darwin-amd64 ssh-web-console-windows-amd64.exe
