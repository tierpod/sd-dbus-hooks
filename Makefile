BINARIES  := bin/sd-dbus-hooks

VERSION ?= 0.1
GITHASH := $(shell git rev-parse --short HEAD)
FULLVER := $(VERSION)-git.$(shell git rev-parse --abbrev-ref HEAD).$(shell git rev-parse --short HEAD)

LDFLAGS := -ldflags "-X main.version=$(FULLVER)"

.PHONY: lint
lint:
	find ./cmd ./pkg -type f -name '*.go' | xargs gofmt -l -e
	go vet ./cmd/... ./pkg/...
	$(GOPATH)/bin/golint ./cmd/... ./pkg/...
	#go test ./cmd/... ./pkg/...

.PHONY: build
build: lint $(BINARIES)

$(BINARIES):
	go build -v $(LDFLAGS) -o $@ cmd/$(notdir $@)/*.go

.PHONY: tag
tag:
	git tag -a -m "Release version" $(VERSION)

.PHONY: clean
clean:
	rm -f bin/*
	rm -f install/*.retry
	rm -f ./pprof

.PHONY: doc
doc:
	godoc -http :6060
