GO           := go
FIRST_GOPATH := $(firstword $(subst :, ,$(GOPATH)))
pkgs         := $(shell $(GO) list ./...)
LDFLAGS      := --ldflags "-X github.com/bpicode/fritzctl/meta.Version=$(FRITZCTL_VERSION)"

all: format build test

format:
	@echo ">> formatting code"
	@$(GO) fmt $(pkgs)

dependencies:
	@echo ">> getting dependencies"
	@$(GO) get -t ./...

build: dependencies
	@echo ">> building project"
	@$(GO) build $(LDFLAGS)

test: build
	@echo ">> testing"
	echo "mode: count" > coverage-all.out
	$(foreach pkg,$(pkgs),\
		go test $(LDFLAGS) -coverprofile=coverage.out -covermode=count $(pkg) || exit 1;\
		tail -n +2 coverage.out >> coverage-all.out;)
		go tool cover -html=coverage-all.out

clean:
	@echo ">> cleaning"
	@$(GO) clean

   	