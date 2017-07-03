GO           := go
FIRST_GOPATH := $(firstword $(subst :, ,$(GOPATH)))
pkgs         := $(shell $(GO) list ./... | grep -v github.com/bpicode/fritzctl/vendor/)
FRITZCTL_VERSION ?= unknown
FRITZCTL_OUTPUT ?= fritzctl
BASH_COMPLETION_OUTPUT ?= "os/completion/fritzctl"
LDFLAGS      := --ldflags "-X github.com/bpicode/fritzctl/config.Version=$(FRITZCTL_VERSION)"
TESTFLAGS    ?=

all: sysinfo format build test completion_bash

sysinfo:
	@echo ">> SYSTEM INFORMATION"
	@echo ">> PLATFORM: $(shell uname -a)"
	@echo ">> GO      : $(shell go version)"

format:
	@echo ">> formatting code"
	@$(GO) fmt $(pkgs)

dependencies:
	@echo ">> getting dependencies"
	@$(GO) get -t -v ./...
	@echo ">> dependencies:"
	@$(GO) list -f '{{join .Deps "\n"}}'

build: dependencies
	@echo ">> building project, version=$(FRITZCTL_VERSION)"
	@$(GO) build -o $(FRITZCTL_OUTPUT) $(LDFLAGS)

test: build
	@echo ">> testing"
	echo "mode: count" > coverage-all.out
	$(foreach pkg,$(pkgs),\
		go test $(LDFLAGS) $(TESTFLAGS) -coverprofile=coverage.out -covermode=atomic $(pkg) || exit 1;\
		tail -n +2 coverage.out >> coverage-all.out;)
		go tool cover -html=coverage-all.out -o coverage-all.html

fasttest: build
	@echo ">> testing, fast mode"
	echo "mode: count" > coverage-all.out
	$(foreach pkg,$(pkgs),\
		go test  $(LDFLAGS) $(TESTFLAGS) -coverprofile=coverage.out $(pkg) || exit 1;\
		tail -n +2 coverage.out >> coverage-all.out;)
		go tool cover -html=coverage-all.out

completion_bash: build
	@echo ">> generating completion script for bash $(BASH_COMPLETION_OUTPUT) using $(FRITZCTL_OUTPUT)"
	$(FRITZCTL_OUTPUT) completion bash > $(BASH_COMPLETION_OUTPUT)

clean:
	@echo ">> cleaning"
	@$(GO) clean

   	
