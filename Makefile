NAME       := sd-dbus-hooks
DESTDIR    := /opt
INSTALLDIR := $(DESTDIR)/$(NAME)

VERSION := $(shell git describe --tags)
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

.PHONY: lint
lint:
	find ./cmd ./pkg -type f -name '*.go' | xargs gofmt -l -e
	go vet ./cmd/... ./pkg/...
	$(GOPATH)/bin/golint ./cmd/... ./pkg/...
	#go test ./cmd/... ./pkg/...

.PHONY: build
build: lint bin/$(NAME)

bin/$(NAME):
	go build -v $(LDFLAGS) -o $@ cmd/$(NAME)/*.go

.PHONY: clean
clean:
	rm -f bin/*
	rm -f install/*.retry
	rm -f ./pprof

.PHONY: doc
doc:
	godoc -http :6060

.PHONY: install
install: $(INSTALLDIR)
	install -m 0755 bin/$(NAME) $(INSTALLDIR)
	install -m 0600 config/config.dist.yaml $(INSTALLDIR)/config.dist.yaml

$(INSTALLDIR) release:
	mkdir -p $@

release/$(NAME)_linux_amd64.tar.gz: release
	make DESTDIR=./tmp install
	tar -cvzf $@ --owner=0 --group=0 -C./tmp $(NAME)
