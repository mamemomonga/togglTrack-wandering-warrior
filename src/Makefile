APPNAME?=togglTrack-wandering-warrior

default: release

bin: bin/$(APPNAME)

bin/$(APPNAME): vendor
	go build -o $@ ./

vendor:
	go mod vendor

release:
	GOOS=darwin GOARCH=amd64 $(MAKE) release-build
	GOOS=darwin GOARCH=arm64 $(MAKE) release-build

release-build:
	go build -o release/$(APPNAME)-$(GOOS)-$(GOARCH)

clean:
	rm -rf vendor release bin
