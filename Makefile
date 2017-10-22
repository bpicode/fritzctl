FIRST_GOPATH := $(firstword $(subst :, ,$(GOPATH)))
PKGS         := $(shell go list ./...)
GOFILES_NOVENDOR := $(shell find . -type f -name '*.go' -not -path "./vendor/*")
FRITZCTL_VERSION ?= unknown
FRITZCTL_OUTPUT ?= fritzctl
BASH_COMPLETION_OUTPUT ?= "os/completion/fritzctl"
MAN_PAGE_OUTPUT ?= "os/man/fritzctl.1"
DEPENDENCIES_GRAPH_OUTPUT ?= "dependencies.png"
LDFLAGS      := --ldflags "-X github.com/bpicode/fritzctl/config.Version=$(FRITZCTL_VERSION)"
TESTFLAGS    ?=

all: sysinfo deps build install test codequality completion_bash man

.PHONY: clean build

define ok
	@tput setaf 6 2>/dev/null || echo -n ""
	@echo " [OK]"
	@tput sgr0 2>/dev/null || echo -n ""
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
	@rm -f ./coverage.out
	@rm -rf ./build/
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

codequality:
	@echo ">> CODE QUALITY"
	@echo -n "     FMT"
	@$(foreach gofile, $(GOFILES_NOVENDOR),\
	        (gofmt -s -l -d -e $(gofile) | tee /dev/stderr) || exit 1;)
	@$(call ok)

	@echo -n "     VET"
	@go vet ./...
	@$(call ok)

	@echo -n "     CYCLO"
	@go get github.com/fzipp/gocyclo
	@$(foreach gofile, $(GOFILES_NOVENDOR),\
			gocyclo -over 15 $(gofile);)
	@$(call ok)
	@echo -n "     LINT"
	@go get github.com/golang/lint/golint
	@$(foreach pkg, $(PKGS),\
			golint -set_exit_status $(pkg);)
	@$(call ok)
	@echo -n "     INEFF"
	@go get github.com/gordonklaus/ineffassign
	@ineffassign .
	@$(call ok)
	@echo -n "     SPELL"
	@go get github.com/client9/misspell/cmd/misspell
	@$(foreach gofile, $(GOFILES_NOVENDOR),\
			misspell --error $(gofile);)
	@$(call ok)

dist_all: dist_linux dist_darwin dist_win

dist_darwin:
	@echo  -n ">> BUILD, darwin/amd64"
	@(GOOS=darwin GOARCH=amd64 go build -o build/distributions/darwin_amd64/fritzctl $(LDFLAGS))
	@$(call ok)

dist_win:
	@echo  -n ">> BUILD, windows/amd64"
	@(GOOS=windows GOARCH=amd64 go build -o build/distributions/windows_amd64/fritzctl.exe $(LDFLAGS))
	@$(call ok)

dist_linux:
	@echo  -n ">> BUILD, linux/amd64"
	@(GOOS=linux GOARCH=amd64 go build -o build/distributions/linux_amd64/usr/bin/fritzctl $(LDFLAGS))
	@$(call ok)

	@echo  -n ">> BUILD, linux/arm"
	@(GOOS=linux GOARCH=arm GOARM=6 go build -o build/distributions/linux_arm/usr/bin/fritzctl $(LDFLAGS))
	@$(call ok)

pkg_all: pkg_linux pkg_darwin pkg_win

pkg_win: dist_win
	@echo  -n ">> PACKAGE, windows/amd64"
	@zip -q build/distributions/fritzctl-$(FRITZCTL_VERSION)-windows-amd64.zip build/distributions/windows_amd64/fritzctl.exe
	@$(call ok)

pkg_darwin: dist_darwin
	@echo  -n ">> PACKAGE, darwin/amd64"
	@zip -q build/distributions/fritzctl-$(FRITZCTL_VERSION)-darwin-amd64.zip build/distributions/darwin_amd64/fritzctl
	@$(call ok)

pkg_linux: dist_linux man completion_bash
	@mkdir -p build/distributions/linux_amd64/usr/bin
	@mkdir -p build/distributions/linux_amd64/etc/fritzctl
	@mkdir -p build/distributions/linux_amd64/etc/bash_completion.d
	@mkdir -p build/distributions/linux_amd64/usr/share/man/man1
	@cp os/completion/fritzctl build/distributions/linux_amd64/etc/bash_completion.d/
	@cp os/config/fritzctl.json build/distributions/linux_amd64/etc/fritzctl/
	@cp os/config/fritz.pem build/distributions/linux_amd64/etc/fritzctl/
	@cp os/man/*.1.gz build/distributions/linux_amd64/usr/share/man/man1/

	@echo ">> PACKAGE, linux/amd64/deb"
	@echo -n "     "
	@$(call mkpkg, amd64, build/distributions/linux_amd64/, build/distributions/, deb)
	@echo ">> PACKAGE, linux/amd64/rpm"
	@echo -n "     "
	@$(call mkpkg, x86_64, build/distributions/linux_amd64/, build/distributions/, rpm)

	@mkdir -p build/distributions/linux_arm/usr/bin
	@mkdir -p build/distributions/linux_arm/etc/fritzctl
	@mkdir -p build/distributions/linux_arm/etc/bash_completion.d
	@mkdir -p build/distributions/linux_arm/usr/share/man/man1
	@cp os/completion/fritzctl build/distributions/linux_arm/etc/bash_completion.d/
	@cp os/config/fritzctl.json build/distributions/linux_arm/etc/fritzctl/
	@cp os/config/fritz.pem build/distributions/linux_arm/etc/fritzctl/
	@cp os/man/*.1.gz build/distributions/linux_arm/usr/share/man/man1/

	@echo ">> PACKAGE, linux/armhf/deb"
	@echo -n "     "
	@$(call mkpkg, armhf, build/distributions/linux_arm/, build/distributions/, deb)
	@echo ">> PACKAGE, linux/arm/rpm"
	@echo -n "     "
	@$(call mkpkg, arm, build/distributions/linux_arm/, build/distributions/, rpm)

define mkpkg
	fpm -f -t $4 -n fritzctl -a $1 -v $(FRITZCTL_VERSION) --log warn --description 'AVM FRITZ!Box client' -m bpicode --vendor bpicode --url https://github.com/bpicode/fritzctl --license MIT --category utils --provides fritzctl --deb-no-default-config-files --config-files etc/fritzctl/fritzctl.json --config-files etc/fritzctl/fritz.pem -p $3 -C $2 -s dir .
endef

sign_deb:
	@echo ">> SIGN, deb packages"
	@echo "     SIGNATURE"
	@dpkg-sig --sign origin -k D0E416CE --g "--no-tty --passphrase=$(DEB_SIGNING_KEY_PASSWORD)" ./build/distributions/*.deb
	@echo "     VERIFY"
	@dpkg-sig --verify ./build/distributions/*.deb

publish_all: publish_deb publish_rpm publish_win

publish_deb:
	@echo ">> PUBLISH, deb packages"

	@$(eval AMD64DEB:=$(shell ls ./build/distributions/fritzctl_*_amd64.deb | xargs -n 1 basename))
	@echo "     UPLOAD -> BINTRAY, $(AMD64DEB)"
	@curl -f -T ./build/distributions/$(AMD64DEB) -ubpicode:$(BINTRAY_API_KEY) -H "X-GPG-PASSPHRASE:$(BINTRAY_SIGN_GPG_PASSPHRASE)" "https://api.bintray.com/content/bpicode/fritzctl_deb/fritzctl/$(FRITZCTL_VERSION)/pool/main/m/fritzctl/$(AMD64DEB);deb_distribution=wheezy,jessie,stretch,sid;deb_component=main;deb_architecture=amd64;publish=1"

	@$(eval ARMDEB:=$(shell ls ./build/distributions/fritzctl_*_armhf.deb | xargs -n 1 basename))
	@echo "     UPLOAD -> BINTRAY, $(AMD64DEB)"
	@curl -f -T ./build/distributions/$(ARMDEB)   -ubpicode:$(BINTRAY_API_KEY) -H "X-GPG-PASSPHRASE:$(BINTRAY_SIGN_GPG_PASSPHRASE)" "https://api.bintray.com/content/bpicode/fritzctl_deb/fritzctl/$(FRITZCTL_VERSION)/pool/main/m/fritzctl/$(ARMDEB);deb_distribution=wheezy,jessie,stretch,sid;deb_component=main;deb_architecture=armhf;publish=1"

	@echo "     CALCULATE METADATA, deb repository"
	@curl -f -X POST -H "X-GPG-PASSPHRASE:$(BINTRAY_SIGN_GPG_PASSPHRASE)" -ubpicode:$(BINTRAY_API_KEY) https://api.bintray.com/calc_metadata/bpicode/fritzctl_deb

publish_rpm:
	@echo ">> PUBLISH, rpm packages"

	@$(eval AMD64RPM:=$(shell ls ./build/distributions/fritzctl-*.x86_64.rpm | xargs -n 1 basename))
	@echo "     UPLOAD -> BINTRAY, $(AMD64RPM)"
	@curl -f -T ./build/distributions/$(AMD64RPM) -ubpicode:$(BINTRAY_API_KEY) -H "X-GPG-PASSPHRASE:$(BINTRAY_SIGN_GPG_PASSPHRASE)" "https://api.bintray.com/content/bpicode/fritzctl_rpm/fritzctl/$(FRITZCTL_VERSION)/$(AMD64RPM);publish=1"

	@$(eval ARMRPM:=$(shell ls ./build/distributions/fritzctl-*.arm.rpm | xargs -n 1 basename))
	@echo "     UPLOAD -> BINTRAY, $(ARMRPM)"
	@curl -f -T ./build/distributions/$(ARMRPM) -ubpicode:$(BINTRAY_API_KEY)  -H "X-GPG-PASSPHRASE:$(BINTRAY_SIGN_GPG_PASSPHRASE)" "https://api.bintray.com/content/bpicode/fritzctl_rpm/fritzctl/$(FRITZCTL_VERSION)/$(ARMRPM);publish=1"

	@echo "     CALCULATE METADATA, rpm repository"
	@curl -f -X POST -H "X-GPG-PASSPHRASE:$(BINTRAY_SIGN_GPG_PASSPHRASE)" -ubpicode:$(BINTRAY_API_KEY) https://api.bintray.com/calc_metadata/bpicode/fritzctl_rpm

publish_win:
	@echo ">> PUBLISH, windows packages"

	@$(eval WINZIP:=$(shell ls ./build/distributions/fritzctl-*.windows-amd64.zip | xargs -n 1 basename))
	@echo "     UPLOAD -> BINTRAY, $(WINZIP)"
	@curl -f -T ./build/distributions/$(WINZIP) -ubpicode:$(BINTRAY_API_KEY) -H "X-GPG-PASSPHRASE:$(BINTRAY_SIGN_GPG_PASSPHRASE)" "https://api.bintray.com/content/bpicode/fritzctl_win/fritzctl/$(FRITZCTL_VERSION)/$(WINZIP);publish=1"


demogif:
	@echo ">> DEMO GIF"
	@go build -o mock/standalone/standalone  mock/standalone/main.go
	@(cd mock/ && standalone/./standalone -httptest.serve=127.0.0.1:8000 & echo $$! > /tmp/TEST_SERVER.PID)
	@sleep 2
	@kill `cat </tmp/TEST_SERVER.PID`

