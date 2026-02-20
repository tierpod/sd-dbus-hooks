NAME       := sd-dbus-hooks
DESTDIR    := /opt
INSTALLDIR := $(DESTDIR)/$(NAME)

GIT_VER    := $(shell git describe --abbrev=7 --always --tags)-$(shell git rev-parse --abbrev-ref HEAD)
LDFLAGS    := -ldflags "-X main.version=$(GIT_VER)-$(shell date +%Y%m%d)"

.PHONY: lint
lint:
	go vet ./app/...
	golangci-lint run ./app/...

.PHONY: test
test:
	go test ./app/...

.PHONY: build
build: lint test bin/$(NAME)

bin/$(NAME):
	go build -v $(LDFLAGS) -o $@ ./app/*.go

.PHONY: clean
clean:
	rm -f bin/*
	rm -f install/*.retry
	rm -f ./pprof
	rm -rf release/*.tar.gz

.PHONY: doc
doc:
	godoc -http :6060

.PHONY: install
install: $(INSTALLDIR)
	install -m 0755 bin/$(NAME) $(INSTALLDIR)
	install -m 0600 config/config.dist.yaml $(INSTALLDIR)/config.dist.yaml

$(INSTALLDIR) release:
	mkdir -p $@

release/$(NAME)_linux_amd64.tar.gz: build release
	make DESTDIR=./tmp install
	tar -cvzf $@ --owner=0 --group=0 -C./tmp $(NAME)
