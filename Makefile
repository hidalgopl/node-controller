BINDIR := $(CURDIR)/bin
LDFLAGS := "-extldflags '-static'"

build:
	GOBIN=$(BINDIR) go install -ldflags $(LDFLAGS) ./...
.PHONY: build

container:
	docker build . -t node-controller:latest
.PHONY: container

test:
	go test ./pkg/node -v
.PHONY: test

lint:
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done
.PHONY: lint