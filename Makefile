APPNAME=togglTrack-wandering-warrior
CONTAINER_IMAGE=toggltrack-wandering-warrior
RELEASE_VERSION=v0.0.1

bin/$(APPNAME): vendor
	go build -o $@ ./

vendor:
	go mod vendor

release-build:
	go build -o release/$(APPNAME)-$(GOOS)-$(GOARCH)

release:
	GOOS=darwin GOARCH=amd64 $(MAKE) release-build
	GOOS=darwin GOARCH=arm64 $(MAKE) release-build

release-publish:
	gh release create $(RELEASE_VERSION) --generate-notes --latest
	gh release upload $(RELEASE_VERSION) release/$(APPNAME)-darwin-amd64
	gh release upload $(RELEASE_VERSION) release/$(APPNAME)-darwin-arm64

build:
	docker build -t $(CONTAINER_IMAGE) .	
	docker run --rm $(CONTAINER_IMAGE) tar cC /g release | tar xv

clean:
	rm -rf bin release

.PHONY: release-build release-publish clean
