GO          ?= GO111MODULE=on go
GO_NOMODULE ?= GO111MODULE=off go

FIRST_GOPATH              := $(firstword $(subst :, ,$(GOPATH)))
PKGS                      := $(shell $(GO) list ./...)
GOFILES_NOVENDOR          := $(shell find . -type f -name '*.go' -not -path "./vendor/*")
FRITZCTL_VERSION          ?= unknown
FRITZCTL_OUTPUT           ?= fritzctl
FRITZCTL_REVISION         := $(shell git rev-parse HEAD)
BASH_COMPLETION_OUTPUT    ?= "os/completion/fritzctl"
MAN_PAGE_OUTPUT           ?= "os/man/fritzctl.1"
COPYRIGHT_OUTPUT          ?= "os/doc/copyright"
BUILDFLAGS                := -ldflags="-s -w -X github.com/bpicode/fritzctl/config.Version=$(FRITZCTL_VERSION) -X github.com/bpicode/fritzctl/config.Revision=$(FRITZCTL_REVISION)" -gcflags="-trimpath=$(GOPATH)" -asmflags="-trimpath=$(GOPATH)"
TESTFLAGS                 ?=

all: sysinfo depverify build install test codequality completion_bash man copyright

.PHONY: clean build man copyright analice

define ok
	@tput setaf 6 2>/dev/null || echo -n ""
	@echo " [OK]"
	@tput sgr0 2>/dev/null || echo -n ""
endef

define lazyinstall
    @which $1 > /dev/null; if [ $$? -ne 0 ]; then \
        $(GO_NOMODULE) get -u $2; \
    fi
endef

sysinfo:
	@echo ">> SYSTEM INFORMATION"
	@echo -n "     PLATFORM: $(shell uname -a)"
	@$(call ok)
	@echo -n "     PWD:    : $(shell pwd)"
	@$(call ok)
	@echo -n "     GO      : $(shell go version)"
	@$(call ok)
	@echo -n "     BUILDFLAGS: $(BUILDFLAGS)"
	@$(call ok)

clean:
	@echo -n ">> CLEAN"
	@$(GO) clean -i
	@rm -f ./os/completion/fritzctl
	@rm -f ./os/man/*.gz
	@rm -f ./os/doc/copyright
	@rm -f ./coverage-all.html
	@rm -f ./coverage-all.out
	@rm -f ./coverage.out
	@rm -rf ./build/
	@rm -f ./fritzctl
	@rm -f ./analice
	@$(call ok)

depverify:
	@echo -n ">> DEPENDENCIES [VERIFY]"
	@$(GO) mod verify 1>/dev/null
	@$(call ok)

build:
	@echo -n ">> BUILD, version = $(FRITZCTL_VERSION)/$(FRITZCTL_REVISION), output = $(FRITZCTL_OUTPUT)"
	@$(GO) build -o $(FRITZCTL_OUTPUT) $(BUILDFLAGS)
	@$(call ok)

install:
	@echo -n ">> INSTALL, version = $(FRITZCTL_VERSION)"
	@$(GO) install $(BUILDFLAGS)
	@$(call ok)

test:
	@echo ">> TEST, \"full-mode\": race detector on"
	@echo "mode: count" > coverage-all.out
	@$(foreach pkg, $(PKGS),\
	    echo -n "     ";\
		$(GO) test -run '(Test|Example)' $(BUILDFLAGS) $(TESTFLAGS) -race -coverprofile=coverage.out -covermode=atomic $(pkg) || exit 1;\
		tail -n +2 coverage.out >> coverage-all.out;)
	@$(GO) tool cover -html=coverage-all.out -o coverage-all.html

fasttest: build
	@echo ">> TEST, \"fast-mode\": race detector off"
	@echo "mode: count" > coverage-all.out
	@$(foreach pkg, $(PKGS),\
	    echo -n "     ";\
		$(GO) test  -run '(Test|Example)' $(BUILDFLAGS) $(TESTFLAGS) -coverprofile=coverage.out $(pkg) || exit 1;\
		tail -n +2 coverage.out >> coverage-all.out;)
	@$(GO) tool cover -html=coverage-all.out -o coverage-all.html

completion_bash:
	@echo -n ">> BASH COMPLETION, output = $(BASH_COMPLETION_OUTPUT)"
	@$(GO) run main.go completion bash > $(BASH_COMPLETION_OUTPUT)
	@$(call ok)

man:
	@echo -n ">> MAN PAGE, output = $(MAN_PAGE_OUTPUT).gz"
	@$(GO) run main.go doc man > $(MAN_PAGE_OUTPUT)
	@gzip --force $(MAN_PAGE_OUTPUT)
	@$(call ok)

analice:
	@echo -n ">> ANALICE"
	@$(GO) build github.com/bpicode/fritzctl/tools/analice
	@$(call ok)

license_compliance: analice
	@echo -n ">> OSS LICENSE COMPLIANCE"
	@$(GO) run github.com/bpicode/fritzctl/tools/analice generate notice $(PKGS) --tests=true --gooses=linux,windows,darwin > NOTICE.tmp
	@diff NOTICE NOTICE.tmp || exit 1
	@rm NOTICE.tmp
	@$(call ok)

copyright: license_compliance
	@echo -n ">> COPYRIGHT, output = $(COPYRIGHT_OUTPUT)"
	@$(GO) run github.com/bpicode/fritzctl/tools/analice generate copyright github.com/bpicode/fritzctl --tests=false --gooses=linux,windows,darwin > $(COPYRIGHT_OUTPUT)
	@$(call ok)

codequality:
	@echo ">> CODE QUALITY"

	@echo -n "     REVIVE"
	@$(call lazyinstall,revive,github.com/mgechev/revive)
	@revive -formatter friendly -exclude vendor/... ./...
	@$(call ok)

	@echo -n "     FMT"
	@$(foreach gofile, $(GOFILES_NOVENDOR),\
	        (gofmt -s -l -d -e $(gofile) | tee /dev/stderr) || exit 1;)
	@$(call ok)

	@echo -n "     VET"
	@$(GO) vet ./...
	@$(call ok)

	@echo -n "     CYCLO"
	@$(call lazyinstall,gocyclo,github.com/fzipp/gocyclo)
	@$(foreach gofile, $(GOFILES_NOVENDOR),\
			gocyclo -over 15 $(gofile);)
	@$(call ok)

	@echo -n "     LINT"
	@$(call lazyinstall,golint,golang.org/x/lint/golint)
	@$(foreach pkg, $(PKGS),\
			golint -set_exit_status $(pkg);)
	@$(call ok)

	@echo -n "     INEFF"
	@$(call lazyinstall,ineffassign,github.com/gordonklaus/ineffassign)
	@ineffassign .
	@$(call ok)

	@echo -n "     SPELL"
	@$(call lazyinstall,misspell,github.com/client9/misspell/cmd/misspell)
	@$(foreach gofile, $(GOFILES_NOVENDOR),\
			misspell --error $(gofile);)
	@$(call ok)

	@echo -n "     STATIC"
	@$(call lazyinstall,staticcheck,honnef.co/go/tools/cmd/staticcheck)
	@staticcheck -checks=all $(PKGS)
	@$(call ok)

	@echo -n "     INTERFACER"
	@$(call lazyinstall,interfacer,mvdan.cc/interfacer)
	@interfacer ./...
	@$(call ok)

	@echo -n "     UNCONVERT"
	@$(call lazyinstall,unconvert,github.com/mdempsky/unconvert)
	@unconvert -v $(PKGS)
	@$(call ok)

dist_all: dist_linux dist_darwin dist_win dist_bsd

define dist
	@echo  -n ">> BUILD, $(1)/$(2) "
	@(GOOS=$(1) GOARCH=$(2) go build -o $(3) $(BUILDFLAGS))
	@cp $(3) build/distributions/fritzctl-$(1)-$(2)$(4)
	@cd build/distributions && shasum -a 256 "fritzctl-$(1)-$(2)$(4)" | tee "fritzctl-$(1)-$(2)$(4).sha256" | cut -b 1-64 | tr -d "\n"
	@$(call ok)
endef

dist_darwin:
	@$(call dist,darwin,amd64,build/distributions/darwin_amd64/fritzctl,"")

dist_win:
	@$(call dist,windows,amd64,build/distributions/windows_amd64/fritzctl.exe,".exe")

dist_linux:
	@$(call dist,linux,amd64,build/distributions/linux_amd64/usr/bin/fritzctl,"")
	@$(call dist,linux,arm,build/distributions/linux_arm/usr/bin/fritzctl,"")

dist_bsd:
	@$(call dist,dragonfly,amd64,build/distributions/dragonfly_amd64/usr/bin/fritzctl,"")
	@$(call dist,freebsd,amd64,build/distributions/freebsd_amd64/usr/bin/fritzctl,"")
	@$(call dist,netbsd,amd64,build/distributions/netbsd_amd64/usr/bin/fritzctl,"")
	@$(call dist,openbsd,amd64,build/distributions/openbsd_amd64/usr/bin/fritzctl,"")

pkg_all: pkg_linux pkg_darwin pkg_win

pkg_win: dist_win
	@echo  -n ">> PACKAGE, windows/amd64"
	@zip -q build/distributions/fritzctl-$(FRITZCTL_VERSION)-windows-amd64.zip build/distributions/windows_amd64/fritzctl.exe
	@$(call ok)

pkg_darwin: dist_darwin
	@echo  -n ">> PACKAGE, darwin/amd64"
	@zip -q build/distributions/fritzctl-$(FRITZCTL_VERSION)-darwin-amd64.zip build/distributions/darwin_amd64/fritzctl
	@$(call ok)

pkg_linux: dist_linux man completion_bash copyright
	@mkdir -p build/distributions/linux_amd64/usr/bin
	@mkdir -p build/distributions/linux_amd64/etc/fritzctl
	@mkdir -p build/distributions/linux_amd64/etc/bash_completion.d
	@mkdir -p build/distributions/linux_amd64/usr/share/man/man1
	@mkdir -p build/distributions/linux_amd64/usr/share/doc/fritzctl
	@cp os/completion/fritzctl build/distributions/linux_amd64/etc/bash_completion.d/
	@cp os/config/config.yml build/distributions/linux_amd64/etc/fritzctl/
	@cp os/config/fritz.pem build/distributions/linux_amd64/etc/fritzctl/
	@cp os/man/*.1.gz build/distributions/linux_amd64/usr/share/man/man1/
	@cp os/doc/copyright build/distributions/linux_amd64/usr/share/doc/fritzctl/

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
	@mkdir -p build/distributions/linux_arm/usr/share/doc/fritzctl
	@cp os/completion/fritzctl build/distributions/linux_arm/etc/bash_completion.d/
	@cp os/config/config.yml build/distributions/linux_arm/etc/fritzctl/
	@cp os/config/fritz.pem build/distributions/linux_arm/etc/fritzctl/
	@cp os/man/*.1.gz build/distributions/linux_arm/usr/share/man/man1/
	@cp os/doc/copyright build/distributions/linux_arm/usr/share/doc/fritzctl/

	@echo ">> PACKAGE, linux/armhf/deb"
	@echo -n "     "
	@$(call mkpkg, armhf, build/distributions/linux_arm/, build/distributions/, deb)
	@echo ">> PACKAGE, linux/arm/rpm"
	@echo -n "     "
	@$(call mkpkg, arm, build/distributions/linux_arm/, build/distributions/, rpm)

define mkpkg
	fpm -f -t $4 -n fritzctl -a $1 -v $(FRITZCTL_VERSION) --log warn --description 'AVM FRITZ!Box client' -m bpicode --vendor bpicode --url https://github.com/bpicode/fritzctl --license MIT --category utils --provides fritzctl --deb-no-default-config-files --config-files etc/fritzctl/config.yml --config-files etc/fritzctl/fritz.pem -p $3 -C $2 -s dir .
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
	@curl -f -T ./build/distributions/$(AMD64DEB) -ubpicode:$(BINTRAY_API_KEY) -H "X-GPG-PASSPHRASE:$(BINTRAY_SIGN_GPG_PASSPHRASE)" "https://api.bintray.com/content/bpicode/fritzctl_deb/fritzctl/$(FRITZCTL_VERSION)/pool/main/m/fritzctl/$(AMD64DEB);deb_distribution=wheezy,jessie,stretch,buster,sid;deb_component=main;deb_architecture=amd64;publish=1"

	@$(eval ARMDEB:=$(shell ls ./build/distributions/fritzctl_*_armhf.deb | xargs -n 1 basename))
	@echo "     UPLOAD -> BINTRAY, $(AMD64DEB)"
	@curl -f -T ./build/distributions/$(ARMDEB)   -ubpicode:$(BINTRAY_API_KEY) -H "X-GPG-PASSPHRASE:$(BINTRAY_SIGN_GPG_PASSPHRASE)" "https://api.bintray.com/content/bpicode/fritzctl_deb/fritzctl/$(FRITZCTL_VERSION)/pool/main/m/fritzctl/$(ARMDEB);deb_distribution=wheezy,jessie,stretch,buster,sid;deb_component=main;deb_architecture=armhf;publish=1"

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

	@$(eval WINZIP:=$(shell ls ./build/distributions/fritzctl-*-windows-amd64.zip | xargs -n 1 basename))
	@echo "     UPLOAD -> BINTRAY, $(WINZIP)"
	@curl -f -T ./build/distributions/$(WINZIP) -ubpicode:$(BINTRAY_API_KEY) -H "X-GPG-PASSPHRASE:$(BINTRAY_SIGN_GPG_PASSPHRASE)" "https://api.bintray.com/content/bpicode/fritzctl_win/fritzctl/$(FRITZCTL_VERSION)/$(WINZIP);publish=1"

demogif:
	@echo ">> DEMO GIF"
	@$(GO) build -o mock/standalone/standalone  mock/standalone/main.go
	@(cd mock/ && standalone/./standalone -httptest.serve=127.0.0.1:8000 & echo $$! > /tmp/TEST_SERVER.PID)
	@sleep 2
	@(cd mock/ && asciinema rec -c '/bin/sh' ../images/fritzctl_demo.json)
	@kill `cat </tmp/TEST_SERVER.PID`
	@docker run --rm -v $(PWD)/images:/data asciinema/asciicast2gif -t monokai fritzctl_demo.json fritzctl_demo.gif

release_github: pkg_all dist_all
	@echo ">> GITHUB RELEASE"
	@$(eval ASSETS:=$(shell find build/ -maxdepth 2 -type f -printf '-a %p\n'))
	@git remote set-url origin https://github.com/bpicode/fritzctl.git
	@hub release create --draft v$(FRITZCTL_VERSION) --message="fritzctl $(FRITZCTL_VERSION)" $(ASSETS)
