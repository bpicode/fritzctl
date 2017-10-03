FIRST_GOPATH := $(firstword $(subst :, ,$(GOPATH)))
PKGS         := $(shell go list ./...)
FRITZCTL_VERSION ?= unknown
FRITZCTL_OUTPUT ?= fritzctl
BASH_COMPLETION_OUTPUT ?= "os/completion/fritzctl"
MAN_PAGE_OUTPUT ?= "os/man/fritzctl.1"
DEPENDENCIES_GRAPH_OUTPUT ?= "dependencies.png"
LDFLAGS      := --ldflags "-X github.com/bpicode/fritzctl/config.Version=$(FRITZCTL_VERSION)"
TESTFLAGS    ?=

all: sysinfo deps build install test completion_bash man

.PHONY: clean build

define ok
	@tput setaf 6
	@echo " [OK]"
	@tput sgr0
endef

sysinfo:
	@echo ">> SYSTEM INFORMATION"
	@echo -n "     PLATFORM: $(shell uname -a)"
	@$(call ok)
	@echo -n "     GO      : $(shell go version)"
	@$(call ok)

clean:
	@echo -n ">> CLEAN"
	@go clean
	@rm -f ./os/completion/fritzctl
	@rm -f ./os/man/*.gz
	@rm -f ./coverage-all.html
	@rm -f ./coverage-all.out
	@$(call ok)

deps:
	@echo -n ">> DEPENDENCIES"
	@go get -u github.com/golang/dep/cmd/dep
	@dep ensure
	@$(call ok)

depprint: deps
	@echo ">> DEPENDENCIES:"
	@dep status

depgraph: deps
	@echo -n ">> DEPENDENCY GRAPH, output = $(DEPENDENCIES_GRAPH_OUTPUT)"
	@dep status -dot | dot -T png -o $(DEPENDENCIES_GRAPH_OUTPUT)
	@$(call ok)

build:
	@echo -n ">> BUILD, version = $(FRITZCTL_VERSION), output = $(FRITZCTL_OUTPUT)"
	@go build -o $(FRITZCTL_OUTPUT) $(LDFLAGS)
	@$(call ok)

install:
	@echo -n ">> INSTALL, version = $(FRITZCTL_VERSION)"
	@go install $(LDFLAGS)
	@$(call ok)

test: build
	@echo ">> TEST, \"full-mode\": race detector on"
	@echo "mode: count" > coverage-all.out
	@$(foreach pkg, $(PKGS),\
	    echo -n "     ";\
		go test $(LDFLAGS) $(TESTFLAGS) -race -coverprofile=coverage.out -covermode=atomic $(pkg) || exit 1;\
		tail -n +2 coverage.out >> coverage-all.out;)
	@go tool cover -html=coverage-all.out -o coverage-all.html

fasttest: build
	@echo ">> TEST, \"fast-mode\": race detector off"
	@echo "mode: count" > coverage-all.out
	@$(foreach pkg, $(PKGS),\
	    echo -n "     ";\
		go test $(LDFLAGS) $(TESTFLAGS) -coverprofile=coverage.out $(pkg) || exit 1;\
		tail -n +2 coverage.out >> coverage-all.out;)
	@go tool cover -html=coverage-all.out

completion_bash: build
	@echo -n ">> BASH COMPLETION, output = $(BASH_COMPLETION_OUTPUT)"
	@$(FRITZCTL_OUTPUT) completion bash > $(BASH_COMPLETION_OUTPUT)
	@$(call ok)

man: build
	@echo -n ">> MAN PAGE, output = $(MAN_PAGE_OUTPUT).gz"
	@$(FRITZCTL_OUTPUT) doc man > $(MAN_PAGE_OUTPUT)
	@gzip --force $(MAN_PAGE_OUTPUT)
	@$(call ok)
