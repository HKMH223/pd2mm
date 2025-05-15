LDFLAGS=-X 'main.buildDate=$(shell date)' -X 'main.gitHash=$(shell git rev-parse --short HEAD)' -X 'main.buildOn=$(shell go version)' -w -s -H=windowsgui -extldflags=-static

GO_BUILD=go build -trimpath -ldflags "$(LDFLAGS)"

BUILD_PATH=./frontend/gui

.PHONY: all fmt mod lint test deadcode syso pd2mm-linux pd2mm-linux-arm64 pd2mm-darwin pd2mm-darwin-arm64 pd2mm-windows clean

all: pd2mm-linux pd2mm-linux-arm64 pd2mm-darwin pd2mm-darwin-arm64 pd2mm-windows

fmt:
	gofumpt -l -w -extra .

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
	windres pd2mm.rc -O coff -o ./frontend/cli/pd2mm.syso
	windres pd2mm.rc -O coff -o ./frontend/gui/pd2mm.syso

pd2mm-linux: fmt
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 $(GO_BUILD) -o pd2mm-linux $(BUILD_PATH)

pd2mm-linux-arm64: fmt
	CGO_ENABLED=1 GOOS=linux GOARCH=arm64 $(GO_BUILD) -o pd2mm-linux-arm64 $(BUILD_PATH)

pd2mm-darwin: fmt
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 $(GO_BUILD) -o pd2mm-darwin $(BUILD_PATH)

pd2mm-darwin-arm64: fmt
	CGO_ENABLED=1 GOOS=linux GOARCH=arm64 $(GO_BUILD) -o pd2mm-darwin-arm64 $(BUILD_PATH)

pd2mm-windows: fmt
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 $(GO_BUILD) -o pd2mm-windows.exe $(BUILD_PATH)

clean:
	rm -f pd2mm-linux pd2mm-linux-arm64 pd2mm-darwin pd2mm-darwin-arm64 pd2mm-windows.exe