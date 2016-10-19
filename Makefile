GO           := go
FIRST_GOPATH := $(firstword $(subst :, ,$(GOPATH)))
pkgs = $(shell $(GO) list ./...)

all: format install test

format:
	@echo ">> formatting code"
	@$(GO) fmt $(pkgs)

dependencies:
	@echo ">> getting dependencies"
	@$(GO) get -t -u ./...

build: dependencies
	@echo ">> building project"
	@$(GO) build

install: dependencies
	@echo ">> installing"
	@$(GO) install ./...

test: dependencies
	@echo ">> testing"
	echo "mode: count" > coverage-all.out
	$(foreach pkg,$(pkgs),\
		go test -coverprofile=coverage.out -covermode=count $(pkg) || exit 1;\
		tail -n +2 coverage.out >> coverage-all.out;)
		go tool cover -html=coverage-all.out

clean:
	@echo ">> cleaning"
	@$(GO) clean

   	