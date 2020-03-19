NAME := bqnotify
VERSION = $(shell gobump show -r)
REVISION := $(shell git rev-parse --short HEAD)

.PHONY: build
build:
	go build -o ./bin/$(NAME)

.PHONY: build
clean:
	rm bin/$(NAME)

.PHONY: lint
lint:
	golint ./...

.PHONY: fmt
fmt:
	goimports -w .

.PHONY: package
package:
	@sh -c "'$(CURDIR)/scripts/package.sh'"

.PHONY: crossbuild
crossbuild:
	goxz -pv=v${VERSION} -build-ldflags="-X main.GitCommit=${REVISION}" \
        -arch=386,amd64 -d=./pkg/dist/v${VERSION} \
        -n ${NAME}
