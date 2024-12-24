VERSION := $(shell git describe --tags)
LDFLAGS += -X "main.BuildTimestamp=$(shell date -u "+%Y-%m-%d %H:%M:%S")"
LDFLAGS += -X "main.Version=$(VERSION)"
LDFLAGS += -X "main.GoVersion=$(shell go version | sed -r 's/go version go(.*)\ .*/\1/')"
GO := GO111MODULE=on CGO_ENABLED=0 go

.PHONY: release
release:
	GOOS=darwin GOARCH=amd64 $(GO) build -ldflags '$(LDFLAGS)' -o bin/darwin-amd64/mjml-dev
	GOOS=darwin GOARCH=arm64 $(GO) build -ldflags '$(LDFLAGS)' -o bin/darwin-arm64/mjml-dev
	GOOS=linux GOARCH=amd64 $(GO) build -ldflags '$(LDFLAGS)' -o bin/linux-amd64/mjml-dev
	GOOS=linux GOARCH=arm64 $(GO) build -ldflags '$(LDFLAGS)' -o bin/linux-arm64/mjml-dev
	GOOS=windows GOARCH=amd64 $(GO) build -ldflags '$(LDFLAGS)' -o bin/windows-amd64/mjml-dev.exe