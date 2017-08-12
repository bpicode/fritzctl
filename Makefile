GO           := go
FIRST_GOPATH := $(firstword $(subst :, ,$(GOPATH)))
pkgs         := $(shell $(GO) list ./...)
FRITZCTL_VERSION ?= unknown
FRITZCTL_OUTPUT ?= fritzctl
BASH_COMPLETION_OUTPUT ?= "os/completion/fritzctl"
MAN_PAGE_OUTPUT ?= "os/man/fritzctl.1"
LDFLAGS      := --ldflags "-X github.com/bpicode/fritzctl/config.Version=$(FRITZCTL_VERSION)"
TESTFLAGS    ?=

all: sysinfo dependencies build test completion_bash man

sysinfo:
	@echo ">> SYSTEM INFORMATION"
	@echo ">> PLATFORM: $(shell uname -a)"
	@echo ">> GO      : $(shell go version)"

dependencies:
	@echo ">> getting dependencies"
	@$(GO) get -u github.com/golang/dep/cmd/dep
	dep ensure
	@echo ">> dependencies:"
	dep status

build:
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

man: build
	@echo ">> generating man page using $(FRITZCTL_OUTPUT)"
	$(FRITZCTL_OUTPUT) doc man > $(MAN_PAGE_OUTPUT)
	gzip --force $(MAN_PAGE_OUTPUT)

clean:
	@echo ">> cleaning"
	@$(GO) clean

   	
