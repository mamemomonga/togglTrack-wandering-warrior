APPNAME  ?= togglTrack-wandering-warrior
VERSION  := $(shell if [ -e version ]; then cat ./version; else cat ../version; fi)
REVISION := $(shell if [ -e revision ]; then cat revision; else git rev-parse --short HEAD; fi)

BUILDINFO_ARGS = -X 'main.version=$(VERSION)' -X 'main.revision=$(REVISION)'

BUILDARGS_DYNAMIC := GO111MODULE=on \
	go build -mod vendor -a -ldflags="-s -w $(BUILDINFO_ARGS)"

BUILDARGS_STATIC := GO111MODULE=on CGO_ENABLED=0 \
	go build -mod vendor -a -tags netgo -installsuffix netgo \
	-ldflags="-s -w $(BUILDINFO_ARGS) -extldflags '-static'"

bin: bin/$(APPNAME)

bin/$(APPNAME): vendor
	$(BUILDARGS_DYNAMIC) -o $@ ./cmd/togglTrack-wandering-warrior

vendor:
	go mod vendor

release:
	GOOS=darwin GOARCH=amd64 $(MAKE) release-build
	GOOS=darwin GOARCH=arm64 $(MAKE) release-build

release-build:
	$(BUILDARGS_STATIC) -o release/$(APPNAME)-$(GOOS)-$(GOARCH) ./cmd/togglTrack-wandering-warrior

clean:
	rm -rf vendor release bin
