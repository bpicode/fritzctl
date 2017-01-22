GO           := go
FIRST_GOPATH := $(firstword $(subst :, ,$(GOPATH)))
pkgs         := $(shell $(GO) list ./...)
FRITZCTL_VERSION ?= unknown
LDFLAGS      := --ldflags "-X github.com/bpicode/fritzctl/meta.Version=$(FRITZCTL_VERSION)"
TESTFLAGS    ?=

all: sysinfo format build test

sysinfo:
	@echo ">> SYSTEM INFORMATION"
	@echo ">> PLATFORM: $(shell uname -a)"
	@echo ">> GO      : $(shell go version)"

format:
	@echo ">> formatting code"
	@$(GO) fmt $(pkgs)

dependencies:
	@echo ">> getting dependencies"
	@$(GO) get -t ./...
	@echo ">> dependencies:"
	@$(GO) list -f '{{join .Deps "\n"}}'

build: dependencies
	@echo ">> building project, version=$(FRITZCTL_VERSION)"
	@$(GO) build $(LDFLAGS)

test: build
	@echo ">> testing"
	echo "mode: count" > coverage-all.out
	$(foreach pkg,$(pkgs),\
		go test $(LDFLAGS) $(TESTFLAGS) -coverprofile=coverage.out -covermode=atomic $(pkg) || exit 1;\
		tail -n +2 coverage.out >> coverage-all.out;)
		go tool cover -html=coverage-all.out

clean:
	@echo ">> cleaning"
	@$(GO) clean

   	
