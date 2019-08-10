BINDIR := $(CURDIR)/bin
LDFLAGS := "-extldflags '-static'"

build:
	GOBIN=$(BINDIR) go install -ldflags $(LDFLAGS) ./...
.PHONY: build

container:
	docker build . -t node-controller:latest
.PHONY: container
