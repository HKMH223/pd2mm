LDFLAGS=-X 'main.buildDate=$(shell date)' -X 'main.gitHash=$(shell git rev-parse --short HEAD)' -X 'main.buildOn=$(shell go version)' -w -s

GO_BUILD=go build -trimpath -ldflags "$(LDFLAGS)"

.PHONY: all fmt mod lint test deadcode syso pd2mm-linux pd2mm-linux-arm64 pd2mm-darwin pd2mm-darwin-arm64 pd2mm-windows clean

all: pd2mm-linux pd2mm-linux-arm64 pd2mm-darwin pd2mm-darwin-arm64 pd2mm-windows

fmt:
	gofumpt -l -w .

mod:
	go get -u
	go mod tidy

lint: fmt
	golangci-lint run --fix

test:
	go test ./...

deadcode:
	deadcode ./...

syso:
	windres pd2mm.rc -O coff -o pd2mm.syso

pd2mm-linux: lint syso
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO_BUILD) -o pd2mm-linux

pd2mm-linux-arm64: lint syso
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GO_BUILD) -o pd2mm-linux-arm64

pd2mm-darwin: lint syso
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO_BUILD) -o pd2mm-darwin

pd2mm-darwin-arm64: lint syso
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GO_BUILD) -o pd2mm-darwin-arm64

pd2mm-windows: lint syso
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO_BUILD) -o pd2mm-windows.exe

clean:
	rm -f pd2mm-linux pd2mm-linux-arm64 pd2mm-darwin pd2mm-darwin-arm64 pd2mm-windows.exe