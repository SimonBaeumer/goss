export GO111MODULE=on

exe = github.com/SimonBaeumer/goss/cmd/goss
pkgs = $(shell ./novendor.sh)
cmd = goss
TRAVIS_TAG ?= "0.0.0"
GO_FILES = $(shell find . \( -path ./vendor -o -name '_test.go' \) -prune -o -name '*.go' -print)

.PHONY: all build install test coverage deps release bench test-int lint gen integration init git-hooks

init: git-hooks

git-hooks:
	$(info INFO: Starting build $@)
	ln -sf ../../.githooks/pre-commit .git/hooks/pre-commit

all: test-all test-all-32

install: release/goss-linux-amd64
	$(info INFO: Starting build $@)
	cp release/$(cmd)-linux-amd64 $(GOPATH)/bin/goss

test:
	$(info INFO: Starting build $@)
	go test $(pkgs)

lint:
	$(info INFO: Starting build $@)
	go test ./...

bench:
	$(info INFO: Starting build $@)
	go test -bench=.

coverage:
	$(info INFO: Starting build $@)
	go test -coverprofile c.out $(pkgs)

release/goss-linux-386: $(GO_FILES)
	$(info INFO: Starting build $@)
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -ldflags "-X main.version=$(TRAVIS_TAG) -s -w" -o release/$(cmd)-linux-386 $(exe)
release/goss-linux-amd64: $(GO_FILES)
	$(info INFO: Starting build $@)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$(TRAVIS_TAG) -s -w" -o release/$(cmd)-linux-amd64 $(exe)
release/goss-linux-arm: $(GO_FILES)
	$(info INFO: Starting build $@)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -ldflags "-X main.version=$(TRAVIS_TAG) -s -w" -o release/$(cmd)-linux-arm $(exe)

release:
	$(MAKE) clean
	$(MAKE) build

build: release/goss-linux-386 release/goss-linux-amd64 release/goss-linux-arm

test-all: lint test integration

test-dgoss:
	$(info INFO: Starting build $@)
	cd extras/dgoss; commander test

deps:
	$(info INFO: Starting build $@)
	go mod vendor

sec:
	$(info INFO: Starting build $@)
	gosec ./...

gen:
	$(info INFO: Starting build $@)
	go generate -tags genny $(pkgs)

clean:
	$(info INFO: Starting build $@)
	rm -rf ./release

integration: build
	$(info INFO: Starting build $@)
	cd integration && commander test

integration-debug:
	$(info INFO: Starting build $@)
	cd integration && commander test --verbose commander.yaml