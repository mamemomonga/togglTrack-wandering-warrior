APPNAME=togglTrack-wandering-warrior
RELEASE_VERSION=v$(shell cat ./version)

CONTAINER_IMAGE=$(shell echo $(APPNAME) | tr A-Z a-z)

release: version
	git rev-parse --short HEAD > revision
	docker build -t $(CONTAINER_IMAGE) .	
	docker run --rm $(CONTAINER_IMAGE) tar cC /g release | tar xv

publish:
	@echo "*** release $(RELEASE_VERSION) ***"
	gh release create $(RELEASE_VERSION) --generate-notes --latest
	gh release upload $(RELEASE_VERSION) release/$(APPNAME)-darwin-amd64
	gh release upload $(RELEASE_VERSION) release/$(APPNAME)-darwin-arm64

clean:
	rm -rf bin release

.PHONY: clean publish

